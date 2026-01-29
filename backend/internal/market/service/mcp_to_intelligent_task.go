package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	iapb "github.com/kymo-mcp/mcpcan/api/market/intelligent_access"
	pb "github.com/kymo-mcp/mcpcan/api/market/mcp_to_intelligent_task"
	"github.com/kymo-mcp/mcpcan/pkg/common"
	"github.com/kymo-mcp/mcpcan/pkg/coze"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/postgres"
	"github.com/kymo-mcp/mcpcan/pkg/dify"
	i18nresp "github.com/kymo-mcp/mcpcan/pkg/i18n"
	"github.com/kymo-mcp/mcpcan/pkg/logger"
	"github.com/kymo-mcp/mcpcan/pkg/n8n"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// McpToIntelligentTaskService struct for mcp to intelligent task service
type McpToIntelligentTaskService struct {
	ctx context.Context
}

// NewMcpToIntelligentTaskService creates a new mcp to intelligent task service
func NewMcpToIntelligentTaskService(ctx context.Context) *McpToIntelligentTaskService {
	go func() {
		// 程序启动的时候加载运行中的任务，启动执行
		tasks, _, err := mysql.McpToIntelligentTaskRepo.FindWithPagination(context.Background(), 1, 1, "", pb.McpToIntelligentTaskStatus_Running.String())
		if err != nil {
			logger.Error(fmt.Sprintf("failed to find mcp to intelligent tasks: %s", err.Error()))
			return
		}
		if len(tasks) > 0 {
			logger.Info(fmt.Sprintf("found mcp to intelligent tasks: %v", tasks))
			ProcessMcpToIntelligentTask(tasks[0].ID)
		}
	}()

	return &McpToIntelligentTaskService{
		ctx: ctx,
	}
}

// CreateHandler creates mcp to intelligent task HTTP handler function
func (s *McpToIntelligentTaskService) CreateHandler(c *gin.Context) {
	var req pb.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}
	// validate request
	if req.Desc == "" || req.IntelligentAccessID == 0 || len(req.InsertIntelligentInfos) == 0 || len(req.McpInstanceIDs) == 0 {
		common.GinError(c, i18nresp.CodeBadRequest, "invalid request")
		return
	}

	// 转换 InsertIntelligentInfos validate request
	var insertInfos model.InsertIntelligentInfos
	for _, info := range req.InsertIntelligentInfos {
		if info.SpaceID == "" || info.UserID == "" || info.SpaceName == "" || info.UserName == "" {
			common.GinError(c, i18nresp.CodeBadRequest, "invalid request")
			return
		}

		setInfo := &model.InsertIntelligentInfo{
			SpaceID:   info.SpaceID,
			UserID:    info.UserID,
			SpaceName: info.SpaceName,
			UserName:  info.UserName,
			Headers:   map[string]*model.HeaderInfo{},
		}
		for key, val := range info.Headers {
			setInfo.Headers[key] = &model.HeaderInfo{
				Token:   val.Token,
				Headers: val.Headers,
			}
		}
		insertInfos = append(insertInfos, setInfo)
	}

	_, total, err := mysql.McpToIntelligentTaskRepo.FindWithPagination(s.ctx, 1, 1, "", pb.McpToIntelligentTaskStatus_Running.String())
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find mcp to intelligent tasks: %s", err.Error()))
		return
	}
	if total > 0 {
		common.GinError(c, i18nresp.CodeBadRequest, "running mcp sync task already exists")
		return
	}

	instances, err := mysql.McpInstanceRepo.FindByInstanceIDs(s.ctx, req.McpInstanceIDs)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find mcp instances: %s", err.Error()))
		return
	}
	if len(instances) != len(req.McpInstanceIDs) {
		common.GinError(c, i18nresp.CodeBadRequest, "invalid request")
		return
	}

	// 获取智能体平台名称
	intelligentAccess, err := mysql.IntelligentAccessRepo.FindByID(s.ctx, req.IntelligentAccessID)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find intelligent access: %s", err.Error()))
		return
	}

	task := &model.McpToIntelligentTask{
		Desc:                   req.Desc,
		IntelligentAccessID:    req.IntelligentAccessID,
		IntelligentAccessName:  intelligentAccess.AccessName,
		InsertIntelligentInfos: insertInfos,
		McpInstanceIDs:         req.McpInstanceIDs,
		Status:                 pb.McpToIntelligentTaskStatus_Running.String(), // 默认状态为运行中
		Domain:                 req.Domain,
		Cookie:                 req.Cookie,
	}

	if err := mysql.McpToIntelligentTaskRepo.Create(s.ctx, task); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to create mcp to intelligent task: %s", err.Error()))
		return
	}

	go func() {
		ProcessMcpToIntelligentTask(task.ID)
	}()

	common.GinSuccess(c, &pb.CreateResponse{
		Task: s.convertToPbTask(task, nil),
	})
}

// CancelHandler cancels mcp to intelligent task HTTP handler function
func (s *McpToIntelligentTaskService) CancelHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}

	// 先查找任务
	task, err := mysql.McpToIntelligentTaskRepo.FindByID(s.ctx, id)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find mcp to intelligent task: %s", err.Error()))
		return
	}

	// 更新状态为取消
	task.Status = pb.McpToIntelligentTaskStatus_Cancel.String()

	if err := mysql.McpToIntelligentTaskRepo.Update(s.ctx, task); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to cancel mcp to intelligent task: %s", err.Error()))
		return
	}

	common.GinSuccess(c, &pb.CancelResponse{})
}

// GetHandler finds mcp to intelligent task by ID HTTP handler function
func (s *McpToIntelligentTaskService) GetHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, fmt.Sprintf("invalid id: %s", err.Error()))
		return
	}

	task, err := mysql.McpToIntelligentTaskRepo.FindByID(s.ctx, id)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find mcp to intelligent task: %s", err.Error()))
		return
	}

	// 查询日志
	logs, err := mysql.McpToIntelligentTaskLogRepo.FindListByTaskID(s.ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find mcp to intelligent task logs: %s", err.Error()))
		return
	}

	common.GinSuccess(c, &pb.GetResponse{
		Task: s.convertToPbTask(task, logs),
	})
}

// ListHandler finds all mcp to intelligent tasks with pagination HTTP handler function
func (s *McpToIntelligentTaskService) ListHandler(c *gin.Context) {
	var req pb.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	tasks, total, err := mysql.McpToIntelligentTaskRepo.FindWithPagination(s.ctx, int(req.Page), int(req.Size), req.Keyword, "")
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find mcp to intelligent tasks: %s", err.Error()))
		return
	}

	var pbTasks []*pb.McpToIntelligentTask
	for _, task := range tasks {
		pbTasks = append(pbTasks, s.convertToPbTask(task, nil))
	}

	common.GinSuccess(c, &pb.ListResponse{
		List:  pbTasks,
		Total: total,
	})
}

// convertToPbTask converts model task to protobuf task
func (s *McpToIntelligentTaskService) convertToPbTask(task *model.McpToIntelligentTask, logs []*model.McpToIntelligentTaskLog) *pb.McpToIntelligentTask {
	// 转换 InsertIntelligentInfos
	var pbInsertInfos []*pb.InsertIntelligentInfo
	for _, info := range task.InsertIntelligentInfos {
		pbInsertInfos = append(pbInsertInfos, &pb.InsertIntelligentInfo{
			SpaceID:   info.SpaceID,
			UserID:    info.UserID,
			SpaceName: info.SpaceName,
			UserName:  info.UserName,
		})
	}

	// 转换 InstallLogs
	var pbInstallLogs []*pb.InstallLog
	if logs != nil {
		var logMap = make(map[string]*pb.InstallLog)
		for _, log := range logs {
			findLog, ok := logMap[log.McpInstanceID]
			if !ok {
				logMap[log.McpInstanceID] = &pb.InstallLog{
					McpInstanceID:         log.McpInstanceID,
					McpInstanceName:       log.McpInstanceName,
					InsertIntelligentLogs: []*pb.InsertIntelligentLog{},
					Status:                log.Status,
					ErrorLog:              log.ErrorLog,
				}
				continue
			}

			if !log.Status && findLog.Status {
				findLog.Status = log.Status
				findLog.ErrorLog = log.ErrorLog
				continue
			}
		}

		for _, log := range logs {
			findLog := logMap[log.McpInstanceID]
			if findLog == nil {
				continue
			}
			findLog.InsertIntelligentLogs = append(findLog.InsertIntelligentLogs, &pb.InsertIntelligentLog{
				InsertIntelligentInfo: &pb.InsertIntelligentInfo{
					SpaceID:   log.SpaceID,
					UserID:    log.UserID,
					SpaceName: log.SpaceName,
					UserName:  log.UserName,
				},
				ErrorLog: log.ErrorLog,
				Status:   log.Status,
			})
		}

		for _, value := range logMap {
			pbInstallLogs = append(pbInstallLogs, value)
		}
		sort.Slice(pbInstallLogs, func(i, j int) bool {
			return pbInstallLogs[i].McpInstanceName < pbInstallLogs[j].McpInstanceName
		})
	}

	// 转换状态
	var status pb.McpToIntelligentTaskStatus
	switch task.Status {
	case pb.McpToIntelligentTaskStatus_Running.String():
		status = pb.McpToIntelligentTaskStatus_Running
	case pb.McpToIntelligentTaskStatus_Success.String():
		status = pb.McpToIntelligentTaskStatus_Success
	case pb.McpToIntelligentTaskStatus_Failed.String():
		status = pb.McpToIntelligentTaskStatus_Failed
	case pb.McpToIntelligentTaskStatus_Cancel.String():
		status = pb.McpToIntelligentTaskStatus_Cancel
	default:
		status = pb.McpToIntelligentTaskStatus_Unknown
	}

	return &pb.McpToIntelligentTask{
		Id:                     task.ID,
		Desc:                   task.Desc,
		IntelligentAccessID:    task.IntelligentAccessID,
		IntelligentAccessName:  task.IntelligentAccessName,
		InsertIntelligentInfos: pbInsertInfos,
		McpInstanceIDs:         task.McpInstanceIDs,
		Status:                 status,
		InstallLogs:            pbInstallLogs,
		CreatedAt:              task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:              task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ProcessMcpToIntelligentTask(id int64) {
	logger.Info("process mcp to intelligent task service", zap.Int64("taskId", id))
	defer logger.Info("process mcp to intelligent task service: finished", zap.Int64("taskId", id))

	// 从数据库中获取任务
	task, err := mysql.McpToIntelligentTaskRepo.FindByID(context.Background(), id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to find mcp to intelligent task: %s", err.Error()), zap.Int64("taskId", id))
		return
	}
	// 从数据库中获取智能体访问信息
	intelligentAccess, err := mysql.IntelligentAccessRepo.FindByID(context.Background(), task.IntelligentAccessID)
	if err != nil {
		if err = mysql.McpToIntelligentTaskRepo.UpdateLogs(context.Background(), id, pb.McpToIntelligentTaskStatus_Failed.String()); err != nil {
			logger.Error(fmt.Sprintf("failed to update mcp to intelligent task logs: %s", err.Error()), zap.Int64("taskId", id))
		}
		logger.Error(fmt.Sprintf("failed to find intelligent access: %s", err.Error()), zap.Int64("IntelligentAccessID", task.IntelligentAccessID))
		return
	}

	// 根据不同平台获取 用户空间 和连接信息
	var userSpaces = []*iapb.UserSpace{}
	var difyConn *sql.DB
	var n8nCookie string
	if intelligentAccess.AccessType == iapb.IntelligentAccessType_COZE.String() {
		userSpaces, err = GetCozeUserSpace(task.Cookie, intelligentAccess)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get coze user space: %s", err.Error()), zap.Int64("IntelligentAccessID", task.IntelligentAccessID))
			if err = mysql.McpToIntelligentTaskRepo.UpdateLogs(context.Background(), id, pb.McpToIntelligentTaskStatus_Failed.String()); err != nil {
				logger.Error(fmt.Sprintf("failed to update mcp to intelligent task logs: %s", err.Error()), zap.Int64("taskId", id))
			}
			return
		}
	} else if intelligentAccess.AccessType == iapb.IntelligentAccessType_N8N.String() {
		n8nCookie, err = n8n.GetCookieFromLogin(intelligentAccess.BaseUrl, intelligentAccess.Username, intelligentAccess.Password)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get n8n cookie: %s", err.Error()), zap.Int64("IntelligentAccessID", task.IntelligentAccessID))
			if err = mysql.McpToIntelligentTaskRepo.UpdateLogs(context.Background(), id, pb.McpToIntelligentTaskStatus_Failed.String()); err != nil {
				logger.Error(fmt.Sprintf("failed to update mcp to intelligent task logs: %s", err.Error()), zap.Int64("taskId", id))
			}
			return
		}
		userSpaces, err = GetN8NUserSpace(intelligentAccess, n8nCookie)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get n8n user space: %s", err.Error()), zap.Int64("IntelligentAccessID", task.IntelligentAccessID))
			if err = mysql.McpToIntelligentTaskRepo.UpdateLogs(context.Background(), id, pb.McpToIntelligentTaskStatus_Failed.String()); err != nil {
				logger.Error(fmt.Sprintf("failed to update mcp to intelligent task logs: %s", err.Error()), zap.Int64("taskId", id))
			}
			return
		}
	} else {
		difyConn, err = BuildTemporaryPostgresConnection(intelligentAccess.DbHost, intelligentAccess.DbPort, intelligentAccess.DbUser, intelligentAccess.DbPassword, intelligentAccess.DbName, false)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to build temporary postgres connection: %s", err.Error()), zap.Int64("IntelligentAccessID", task.IntelligentAccessID))
			if err = mysql.McpToIntelligentTaskRepo.UpdateLogs(context.Background(), id, pb.McpToIntelligentTaskStatus_Failed.String()); err != nil {
				logger.Error(fmt.Sprintf("failed to update mcp to intelligent task logs: %s", err.Error()), zap.Int64("taskId", id))
			}
			return
		}
		defer difyConn.Close()

		userSpaces, err = GetDifyUserSpace(difyConn)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get dify user space: %s", err.Error()), zap.Int64("IntelligentAccessID", task.IntelligentAccessID))
			if err = mysql.McpToIntelligentTaskRepo.UpdateLogs(context.Background(), id, pb.McpToIntelligentTaskStatus_Failed.String()); err != nil {
				logger.Error(fmt.Sprintf("failed to update mcp to intelligent task logs: %s", err.Error()), zap.Int64("taskId", id))
			}
			return
		}
	}
	// 最后状态，如果中间有失败的则将状态设置为失败
	var lastStatus = pb.McpToIntelligentTaskStatus_Success.String()

	// 并发添加，增加速度
	var wg sync.WaitGroup
	// 并发执行数量
	concurrency := 10
	sem := make(chan struct{}, concurrency)

	for _, instanceID := range task.McpInstanceIDs {
		// 获取实例信息
		mcpInstance, err := mysql.McpInstanceRepo.FindByInstanceID(context.Background(), instanceID)
		if err != nil {
			log := model.McpToIntelligentTaskLog{
				TaskID:          task.ID,
				McpInstanceID:   instanceID,
				McpInstanceName: "unknown",
				Status:          false,
				ErrorLog:        fmt.Sprintf("failed to find mcpInstance: %s", err.Error()),

				IntelligentAccessName: task.IntelligentAccessName,
				IntelligentAccessID:   task.IntelligentAccessID,

				SpaceID:   task.InsertIntelligentInfos[0].SpaceID,
				UserID:    task.InsertIntelligentInfos[0].UserID,
				SpaceName: task.InsertIntelligentInfos[0].SpaceName,
				UserName:  task.InsertIntelligentInfos[0].UserName,
			}
			err := mysql.McpToIntelligentTaskLogRepo.Create(context.Background(), &log)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to create mcp to intelligent task log: %s", err.Error()), zap.Int64("taskId", id))
			}

			lastStatus = pb.McpToIntelligentTaskStatus_Failed.String()
			continue
		}

		// 获取实例令牌
		tokens, err := mysql.McpTokenRepo.ListByInstanceID(context.Background(), mcpInstance.InstanceID)
		if err != nil {
			log := model.McpToIntelligentTaskLog{
				TaskID:          task.ID,
				McpInstanceID:   instanceID,
				McpInstanceName: mcpInstance.InstanceName,
				Status:          false,
				ErrorLog:        fmt.Sprintf("failed to find mcpInstance tokens: %s", err.Error()),

				IntelligentAccessName: task.IntelligentAccessName,
				IntelligentAccessID:   task.IntelligentAccessID,

				SpaceID:   task.InsertIntelligentInfos[0].SpaceID,
				UserID:    task.InsertIntelligentInfos[0].UserID,
				SpaceName: task.InsertIntelligentInfos[0].SpaceName,
				UserName:  task.InsertIntelligentInfos[0].UserName,
			}
			err := mysql.McpToIntelligentTaskLogRepo.Create(context.Background(), &log)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to create mcp to intelligent task log: %s", err.Error()), zap.Int64("taskId", id))
			}
			lastStatus = pb.McpToIntelligentTaskStatus_Failed.String()
			continue
		}

		// 遍历插入信息列表
		for _, insertInfo := range task.InsertIntelligentInfos {
			wg.Add(1)
			sem <- struct{}{}
			go func(insertInfo *model.InsertIntelligentInfo) {
				defer func() {
					<-sem
					wg.Done()
				}()

				// 获取任务最新状态，如果是取消则停止
				searchTask, err := mysql.McpToIntelligentTaskRepo.FindByID(context.Background(), id)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to find mcp to intelligent task: %s", err.Error()), zap.Int64("taskId", id))
				} else {
					// 监控任务是否取消，如果取消则直接跳出循环
					if searchTask.Status == pb.McpToIntelligentTaskStatus_Cancel.String() {
						return
					}
				}

				// 给实例先创建对应的 token
				if insertInfo.Headers[instanceID] != nil {
					err = createOrUpdateInstanceToken(instanceID, insertInfo, task, intelligentAccess, tokens)
					if err != nil {
						log := model.McpToIntelligentTaskLog{
							TaskID:          task.ID,
							McpInstanceID:   instanceID,
							McpInstanceName: mcpInstance.InstanceName,
							Status:          false,
							ErrorLog:        err.Error(),

							IntelligentAccessName: task.IntelligentAccessName,
							IntelligentAccessID:   task.IntelligentAccessID,

							SpaceID:   insertInfo.SpaceID,
							UserID:    insertInfo.UserID,
							SpaceName: insertInfo.SpaceName,
							UserName:  insertInfo.UserName,
						}
						err := mysql.McpToIntelligentTaskLogRepo.Create(context.Background(), &log)
						if err != nil {
							logger.Error(fmt.Sprintf("failed to create mcp to intelligent task log: %s", err.Error()), zap.Int64("taskId", id))
						}
						lastStatus = pb.McpToIntelligentTaskStatus_Failed.String()
						return
					}
				}

				// 根据不同的平台创建对应的mcp插件
				if intelligentAccess.AccessType == iapb.IntelligentAccessType_COZE.String() {
					err = createCozeTools(task.Domain, insertInfo, mcpInstance, userSpaces, task.Cookie)
				} else if intelligentAccess.AccessType == iapb.IntelligentAccessType_N8N.String() {
					// 执行创建 n8n tools
					err = createN8NTools(task.Domain, insertInfo, mcpInstance, userSpaces, intelligentAccess, n8nCookie)
				} else {
					// 执行创建 dify tools
					err = createDifyTools(task.Domain, insertInfo, mcpInstance, userSpaces, difyConn)
				}
				if err != nil {
					log := model.McpToIntelligentTaskLog{
						TaskID:          task.ID,
						McpInstanceID:   instanceID,
						McpInstanceName: mcpInstance.InstanceName,
						Status:          false,
						ErrorLog:        err.Error(),

						IntelligentAccessName: task.IntelligentAccessName,
						IntelligentAccessID:   task.IntelligentAccessID,

						SpaceID:   insertInfo.SpaceID,
						UserID:    insertInfo.UserID,
						SpaceName: insertInfo.SpaceName,
						UserName:  insertInfo.UserName,
					}
					err := mysql.McpToIntelligentTaskLogRepo.Create(context.Background(), &log)
					if err != nil {
						logger.Error(fmt.Sprintf("failed to create mcp to intelligent task log: %s", err.Error()), zap.Int64("taskId", id))
					}
					lastStatus = pb.McpToIntelligentTaskStatus_Failed.String()
					return
				}

				// 添加成功的日志
				log := model.McpToIntelligentTaskLog{
					TaskID:          task.ID,
					McpInstanceID:   instanceID,
					McpInstanceName: mcpInstance.InstanceName,
					Status:          true,
					ErrorLog:        "",

					IntelligentAccessName: task.IntelligentAccessName,
					IntelligentAccessID:   task.IntelligentAccessID,

					SpaceID:   insertInfo.SpaceID,
					UserID:    insertInfo.UserID,
					SpaceName: insertInfo.SpaceName,
					UserName:  insertInfo.UserName,
				}
				err = mysql.McpToIntelligentTaskLogRepo.Create(context.Background(), &log)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to create mcp to intelligent task log: %s", err.Error()), zap.Int64("taskId", id))
				}
			}(insertInfo)
		}
		wg.Wait()
	}

	// 最后更新整个任务的状态
	if err = mysql.McpToIntelligentTaskRepo.UpdateLogs(context.Background(), id, lastStatus); err != nil {
		logger.Error(fmt.Sprintf("failed to update mcp to intelligent task logs: %s", err.Error()), zap.Int64("taskId", id))
	}
}

func createN8NTools(domain string, insertInfo *model.InsertIntelligentInfo, mcpInstance *model.McpInstance, userSpaces []*iapb.UserSpace, access *model.IntelligentAccess, cookie string) error {
	var findUserSpace = findUserSpace(userSpaces, insertInfo.UserID, insertInfo.SpaceID)
	if findUserSpace == nil {
		return fmt.Errorf("failed to find n8n user space")
	}

	mcpServerUrl := fmt.Sprintf("%s%s", domain, mcpInstance.PublicProxyPath)

	var findN8nCredential *n8n.CredentialItem
	n8nCredentials, err := n8n.GetCredentials(access.BaseUrl, cookie, insertInfo.SpaceID)
	if err != nil {
		return fmt.Errorf("failed to get n8n credentials: %s", err.Error())
	}
	for _, n8nCredential := range n8nCredentials.Data {
		if n8nCredential.Data["httpStreamUrl"] == mcpServerUrl {
			findN8nCredential = &n8nCredential
			break
		}
	}

	token := ""
	insertHeader := insertInfo.Headers[mcpInstance.InstanceID]
	if insertHeader != nil {
		token = insertHeader.Token
	}

	if findN8nCredential == nil {
		_, err = n8n.CreateMCPCredential(access.BaseUrl, cookie, mcpInstance.InstanceName, mcpServerUrl, fmt.Sprintf("Authorization=%s", token), findUserSpace.TenantID)
		if err != nil {
			return fmt.Errorf("failed to create n8n mcp credential: %s", err.Error())
		}
	} else {
		_, err = n8n.UpdateCredential(access.BaseUrl, cookie, findN8nCredential.ID, mcpInstance.InstanceName, mcpServerUrl, fmt.Sprintf("Authorization=%s", token))
		if err != nil {
			return fmt.Errorf("failed to update n8n mcp credential: %s", err.Error())
		}
	}
	return nil
}

func createCozeTools(domain string, insertInfo *model.InsertIntelligentInfo, mcpInstance *model.McpInstance, userSpaces []*iapb.UserSpace, cookie string) error {
	// 获取对应的用户空间
	var findUserSpace = findUserSpace(userSpaces, insertInfo.UserID, insertInfo.SpaceID)
	if findUserSpace == nil {
		return fmt.Errorf("failed to find coze user space")
	}

	// 获取 coze 插件列表
	pluginList, err := coze.GetPluginList(cookie, insertInfo.SpaceID)
	if err != nil {
		return fmt.Errorf("failed to get coze plugin list: %s", err.Error())
	}

	// 查询是否已经存在对应的插件
	var pluginInfo *coze.GetPluginInfoResponse
	var pluginID string
	mcpServerUrl := fmt.Sprintf("%s%s", domain, mcpInstance.PublicProxyPath)
	for _, plugin := range pluginList {
		info, err := coze.GetPluginInfo(cookie, plugin.ResID)
		if err != nil {
			continue
		}
		if info.MetaInfo.URL == mcpServerUrl {
			pluginInfo = info
			pluginID = plugin.ResID
			break
		}
	}

	// 描述和简介截断，防止超长
	desc := mcpInstance.Notes
	driefIntro := ""
	if desc == "" {
		desc = mcpInstance.InstanceName
	}
	if len(desc) >= 50 {
		driefIntro = desc[:50]
	}

	// 创建
	if pluginInfo == nil {
		token := ""
		insertHeader := insertInfo.Headers[mcpInstance.InstanceID]
		if insertHeader != nil {
			token = insertHeader.Token
		}

		name := mcpInstance.InstanceName
		if len(name) >= 30 {
			name = mcpInstance.InstanceName[:30]
		}

		resp, err := coze.RegisterPlugin(&coze.RegisterPluginRequest{
			Name: name,
			Desc: desc,
			URL:  mcpServerUrl,
			Icon: coze.IconRequest{
				URI: "plugin_icon/default_icon.png",
			},
			AuthType:  0,
			OauthInfo: "{}",
			SpaceID:   insertInfo.SpaceID,
			CommonParams: map[string][]coze.ParamItem{
				"1": {},
				"2": {},
				"3": {},
				"4": {
					{
						Name:  "Authorization",
						Value: token,
					},
				},
			},
			IdeCodeRuntime: "1",
			PluginType:     11,
			BriefIntro:     driefIntro,
		}, cookie)
		if err != nil {
			return fmt.Errorf("failed to register plugin: %s", err.Error())
		}
		pluginID = resp.PluginID
	} else {
		// 更新
		_, err := coze.UpdatePlugin(&coze.UpdatePluginRequest{
			PluginID: pluginID,
			Name:     mcpInstance.InstanceName,
			Desc:     desc,
			URL:      pluginInfo.MetaInfo.URL,
			Icon: coze.IconRequest{
				URI: pluginInfo.MetaInfo.Icon.URI,
				URL: pluginInfo.MetaInfo.Icon.URL,
			},
			AuthType:     pluginInfo.MetaInfo.AuthType[0],
			OAuthInfo:    pluginInfo.MetaInfo.OauthInfo,
			CommonParams: pluginInfo.MetaInfo.CommonParams,
			EditVersion:  pluginInfo.EditVersion,
			PluginType:   11,
			BriefIntro:   driefIntro,
		}, cookie)
		if err != nil {
			return fmt.Errorf("failed to update plugin: %s", err.Error())
		}
	}

	// 更新插件工具列表
	_, err = coze.RefreshToolList(cookie, pluginID, insertInfo.SpaceID)
	if err != nil {
		return fmt.Errorf("failed to refresh coze tool list: %s", err.Error())
	}
	// 发布插件
	time := time.Now().Format("20060102150405")
	_, err = coze.PublishPlugin(cookie, &coze.PublishRequest{
		PluginID:      pluginID,
		PrivacyStatus: false,
		PrivacyInfo:   "",
		VersionName:   fmt.Sprintf("v0.0.%s", time),
		VersionDesc:   "mcpcan publish",
	})
	if err != nil {
		return fmt.Errorf("failed to publish plugin: %s", err.Error())
	}
	return nil
}

func findUserSpace(userSpaces []*iapb.UserSpace, userID string, spaceID string) *iapb.UserSpace {
	var findUserSpace *iapb.UserSpace
	for _, userSpace := range userSpaces {
		if userSpace.UserID == userID && userSpace.TenantID == spaceID {
			findUserSpace = userSpace
			break
		}
	}
	return findUserSpace
}

func createOrUpdateInstanceToken(instanceID string, insertInfo *model.InsertIntelligentInfo, task *model.McpToIntelligentTask, intelligentAccess *model.IntelligentAccess, tokens []*model.McpToken) error {
	headerJson, _ := json.Marshal(insertInfo.Headers[instanceID].Headers)
	// 给 instance 创建或者更新对应的 token
	usageToken := FindToken(tokens, insertInfo, task.IntelligentAccessID)
	if usageToken == nil {
		labels := []string{
			fmt.Sprintf("user_id=%v", insertInfo.UserID),
			fmt.Sprintf("user_name=%v", insertInfo.UserName),
			fmt.Sprintf("space_id=%v", insertInfo.SpaceID),
			fmt.Sprintf("space_name=%v", insertInfo.SpaceName),
			fmt.Sprintf("intelligent_access_id=%v", task.IntelligentAccessID),
			fmt.Sprintf("intelligent_access_name=%v", task.IntelligentAccessName),
			fmt.Sprintf("intelligent_access_type=%v", intelligentAccess.AccessType),
		}

		labelsJson, _ := json.Marshal(labels)
		usageToken = &model.McpToken{
			InstanceID: instanceID,
			Token:      insertInfo.Headers[instanceID].Token,
			Enabled:    true,
			Headers:    headerJson,
			Usages:     labelsJson,
			ExpireAt:   0,
			PublishAt:  time.Now().UnixMilli(),
		}
		err := mysql.McpTokenRepo.Create(context.Background(), usageToken)
		if err != nil {
			return fmt.Errorf("failed to create token: %s", err.Error())
		}
	} else {
		usageToken.Headers = headerJson
		err := mysql.McpTokenRepo.Update(context.Background(), usageToken)
		if err != nil {
			return fmt.Errorf("failed to update token: %s", err.Error())
		}
	}
	return nil
}

func createDifyTools(domain string, insertInfo *model.InsertIntelligentInfo, mcpInstance *model.McpInstance, userSpaces []*iapb.UserSpace, difyConn *sql.DB) error {
	// 获取对应的用户空间
	var findUserSpace *iapb.UserSpace
	for _, userSpace := range userSpaces {
		if userSpace.UserID == insertInfo.UserID && userSpace.TenantID == insertInfo.SpaceID {
			findUserSpace = userSpace
			break
		}
	}
	if findUserSpace == nil {
		return fmt.Errorf("failed to find dify user space")
	}

	// 获取该插入信息的对应的实例要传递的 header
	listToolsHeaders := map[string]string{}
	gatewayHeader := map[string]string{}

	insertHeader := insertInfo.Headers[mcpInstance.InstanceID]
	if insertHeader != nil {
		listToolsHeaders = insertHeader.Headers

		token, err := dify.EncryptToken(findUserSpace.EncryptPublicKey, insertHeader.Token)
		if err != nil {
			return fmt.Errorf("failed to encrypt dify token: %s", err.Error())
		}
		gatewayHeader["Authorization"] = token
	}

	// 给该 mcp 实例创建对应的 http client
	mcpClient, err := client.NewStreamableHttpClient(
		mcpInstance.ContainerServiceURL,
		transport.WithHTTPHeaders(listToolsHeaders),
	)
	if err != nil {
		return fmt.Errorf("failed to call mcp failed: %s", err.Error())
	}
	_, err = mcpClient.Initialize(context.Background(), mcp.InitializeRequest{})
	if err != nil {
		return fmt.Errorf("failed to call mcp failed, init mcp failed: %s", err.Error())
	}
	// 调用 mcp 服务的 list tools 接口
	tools, err := mcpClient.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		return fmt.Errorf("failed to call mcp, list tools failed: %s", err.Error())
	}
	toolsJson, err := json.Marshal(tools.Tools)
	if err != nil {
		return fmt.Errorf("failed to marshal tools failed: %s", err.Error())
	}

	mcpServerUrl := fmt.Sprintf("%s%s", domain, mcpInstance.PublicProxyPath)
	serverURL, err := dify.EncryptToken(findUserSpace.EncryptPublicKey, mcpServerUrl)
	if err != nil {
		return fmt.Errorf("encrypt token failed: %s", err.Error())
	}
	mcpServerUrlHash := computeSHA256Hash(mcpServerUrl)
	provider, err := postgres.GetToolMcpProvider(difyConn, mcpServerUrlHash, insertInfo.SpaceID)
	if err != nil {
		return fmt.Errorf("get dify provider failed: %s", err.Error())
	}
	// 创建
	if provider == nil {
		gatewayHeaderJson, _ := json.Marshal(gatewayHeader)
		provider = &postgres.ToolMcpProvider{
			Name:                 mcpInstance.InstanceName,
			ServerIdentifier:     mcpInstance.InstanceID[:23],
			ServerURL:            serverURL,
			ServerURLHash:        mcpServerUrlHash,
			Icon:                 "",
			TenantId:             insertInfo.SpaceID,
			UserId:               insertInfo.UserID,
			EncryptedCredentials: "{}",
			Authed:               true,
			Tools:                string(toolsJson),
			Timeout:              30,
			SseReadTimeout:       300,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			EncryptedHeaders:     string(gatewayHeaderJson),
		}
		err = postgres.CreateToolMcpProvider(difyConn, provider)
		if err != nil {
			return fmt.Errorf("failed to create provider failed: %s", err.Error())
		}
	} else {
		provider.Name = mcpInstance.InstanceName
		provider.Tools = string(toolsJson)
		provider.UpdatedAt = time.Now()
		err := postgres.UpdateToolMcpProvider(difyConn, provider)
		if err != nil {
			return fmt.Errorf("failed to update provider failed: %s", err.Error())
		}
	}
	return nil
}

func computeSHA256Hash(input string) string {
	// 创建 SHA256 哈希对象
	hasher := sha256.New()

	// 写入数据
	hasher.Write([]byte(input))

	// 计算哈希并转换为十六进制字符串
	return hex.EncodeToString(hasher.Sum(nil))
}

func FindToken(tokens []*model.McpToken, info *model.InsertIntelligentInfo, intelligentAccessID int64) *model.McpToken {
	for _, token := range tokens {
		if len(token.Usages) == 0 {
			continue
		}

		var usages []string
		if len(token.Usages) > 0 {
			_ = json.Unmarshal(token.Usages, &usages)
		}

		matchUser := false
		for _, usage := range usages {
			if usage == fmt.Sprintf("user_id=%s", info.UserID) {
				matchUser = true
			}
		}

		matchSpace := false
		for _, usage := range usages {
			if usage == fmt.Sprintf("space_id=%s", info.SpaceID) {
				matchSpace = true
			}
		}

		matchIntelligentAccess := false
		for _, usage := range usages {
			if usage == fmt.Sprintf("intelligent_access_id=%v", intelligentAccessID) {
				matchIntelligentAccess = true
			}
		}

		if matchSpace && matchUser && matchIntelligentAccess {
			return token
		}
	}
	return nil
}
