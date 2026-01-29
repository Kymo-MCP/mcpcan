package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pb "github.com/kymo-mcp/mcpcan/api/market/intelligent_access"
	"github.com/kymo-mcp/mcpcan/internal/market/config"
	"github.com/kymo-mcp/mcpcan/pkg/common"
	"github.com/kymo-mcp/mcpcan/pkg/coze"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/postgres"
	i18nresp "github.com/kymo-mcp/mcpcan/pkg/i18n"
	"github.com/kymo-mcp/mcpcan/pkg/n8n"
	_ "github.com/lib/pq"
)

const N8NPluginName = "n8n-nodes-mcp"

// IntelligentAccessService struct for intelligent access service
type IntelligentAccessService struct {
	ctx context.Context
}

// NewIntelligentAccessService creates a new intelligent access service
func NewIntelligentAccessService(ctx context.Context) *IntelligentAccessService {
	return &IntelligentAccessService{
		ctx: ctx,
	}
}

// CreateHandler creates intelligent access HTTP handler function
func (s *IntelligentAccessService) CreateHandler(c *gin.Context) {
	if config.GetConfig().RunMode == common.RunModeKymo {
		common.GinError(c, i18nresp.CodeForbidden, "cannot create intelligent access when running in kimo mode")
		return
	}

	var req pb.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}
	if pb.IntelligentAccessType_value[req.AccessType] <= 0 {
		common.GinError(c, i18nresp.CodeBadRequest, "invalid access type")
		return
	}

	// validate request
	if req.AccessType == pb.IntelligentAccessType_DifyEnterprise.String() || req.AccessType == pb.IntelligentAccessType_Dify.String() || req.AccessType == pb.IntelligentAccessType_QAgent.String() {
		if req.AccessName == "" || req.AccessType == "" || req.DbHost == "" || req.DbPort == 0 || req.DbUser == "" || req.DbPassword == "" || req.DbName == "" {
			common.GinError(c, i18nresp.CodeBadRequest, "invalid request")
			return
		}
	}

	if req.AccessType == pb.IntelligentAccessType_COZE.String() {
		if pb.SubType_value[req.SubType] <= 0 {
			common.GinError(c, i18nresp.CodeBadRequest, "invalid sub type")
			return
		}
		if pb.SubType_value[req.SubType] == int32(pb.SubType_Team) {
			if req.EnterpriseID == "" {
				common.GinError(c, i18nresp.CodeBadRequest, "invalid enterprise id")
				return
			}
		} else if pb.SubType_value[req.SubType] == int32(pb.SubType_Person) {
			if req.CozeUserID == "" {
				common.GinError(c, i18nresp.CodeBadRequest, "invalid coze user id")
				return
			}
		}
	}
	if req.AccessType == pb.IntelligentAccessType_N8N.String() {
		if req.BaseUrl == "" || req.Username == "" || req.Password == "" {
			common.GinError(c, i18nresp.CodeBadRequest, "invalid request")
			return
		}
	}

	intelligentAccess := &model.IntelligentAccess{
		AccessName:   req.AccessName,
		AccessType:   req.AccessType,
		DbHost:       req.DbHost,
		DbPort:       int(req.DbPort),
		DbUser:       req.DbUser,
		DbPassword:   req.DbPassword,
		DbName:       req.DbName,
		SubType:      req.SubType,
		EnterpriseID: req.EnterpriseID,
		BaseUrl:      req.BaseUrl,
		Username:     req.Username,
		Password:     req.Password,
		CozeUserID:   req.CozeUserID,
	}

	if err := mysql.IntelligentAccessRepo.Create(s.ctx, intelligentAccess); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to create intelligent access: %s", err.Error()))
		return
	}

	common.GinSuccess(c, &pb.CreateResponse{
		IntelligentAccess: CoverDbAccessToPn(intelligentAccess),
	})
}

// UpdateHandler updates intelligent access HTTP handler function
func (s *IntelligentAccessService) UpdateHandler(c *gin.Context) {
	if config.GetConfig().RunMode == common.RunModeKymo {
		common.GinError(c, i18nresp.CodeForbidden, "cannot update intelligent access when running in kimo mode")
		return
	}

	var req pb.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}

	dbIntelligentAccess, err := mysql.IntelligentAccessRepo.FindByID(s.ctx, req.AccessID)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find intelligent access: %s", err.Error()))
		return
	}
	if dbIntelligentAccess.AccessType == pb.IntelligentAccessType_QAgent.String() || dbIntelligentAccess.AccessType == pb.IntelligentAccessType_DifyEnterprise.String() || dbIntelligentAccess.AccessType == pb.IntelligentAccessType_Dify.String() {
		if req.AccessID == 0 || req.AccessName == "" || req.DbHost == "" || req.DbPort == 0 || req.DbUser == "" || req.DbPassword == "" || req.DbName == "" {
			common.GinError(c, i18nresp.CodeBadRequest, "invalid request")
			return
		}

		dbIntelligentAccess.AccessName = req.AccessName
		dbIntelligentAccess.DbHost = req.DbHost
		dbIntelligentAccess.DbName = req.DbName
		dbIntelligentAccess.DbPassword = req.DbPassword
		dbIntelligentAccess.DbPort = int(req.DbPort)
		dbIntelligentAccess.DbUser = req.DbUser
	}
	if dbIntelligentAccess.AccessType == pb.IntelligentAccessType_COZE.String() {
		if pb.SubType_value[req.SubType] <= 0 {
			common.GinError(c, i18nresp.CodeBadRequest, "invalid sub type")
			return
		}
		if pb.SubType_value[req.SubType] == int32(pb.SubType_Team) {
			if req.EnterpriseID == "" {
				common.GinError(c, i18nresp.CodeBadRequest, "invalid enterprise id")
				return
			}

			dbIntelligentAccess.EnterpriseID = req.EnterpriseID
		} else if pb.SubType_value[req.SubType] == int32(pb.SubType_Person) {
			if req.CozeUserID == "" {
				common.GinError(c, i18nresp.CodeBadRequest, "invalid coze user id")
				return
			}

			dbIntelligentAccess.CozeUserID = req.CozeUserID
		}

		dbIntelligentAccess.SubType = req.SubType
		dbIntelligentAccess.AccessName = req.AccessName
	}
	if dbIntelligentAccess.AccessType == pb.IntelligentAccessType_N8N.String() {
		if req.BaseUrl == "" || req.Username == "" || req.Password == "" {
			common.GinError(c, i18nresp.CodeBadRequest, "invalid request")
			return
		}

		dbIntelligentAccess.BaseUrl = req.BaseUrl
		dbIntelligentAccess.Username = req.Username
		dbIntelligentAccess.Password = req.Password
	}

	if err := mysql.IntelligentAccessRepo.Update(s.ctx, dbIntelligentAccess); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to update intelligent access: %s", err.Error()))
		return
	}

	common.GinSuccess(c, &pb.UpdateResponse{
		IntelligentAccess: CoverDbAccessToPn(dbIntelligentAccess),
	})
}

// DeleteHandler deletes intelligent access HTTP handler function
func (s *IntelligentAccessService) DeleteHandler(c *gin.Context) {
	if config.GetConfig().RunMode == common.RunModeKymo {
		common.GinError(c, i18nresp.CodeForbidden, "cannot delete intelligent access when running in kimo mode")
		return
	}

	var req pb.DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}

	if err := mysql.IntelligentAccessRepo.Delete(s.ctx, req.AccessID); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to delete intelligent access: %s", err.Error()))
		return
	}

	common.GinSuccess(c, &pb.DeleteResponse{})
}

// GetHandler finds intelligent access by ID HTTP handler function
func (s *IntelligentAccessService) GetHandler(c *gin.Context) {
	var req pb.GetRequest
	if err := c.ShouldBindUri(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}

	intelligentAccess, err := mysql.IntelligentAccessRepo.FindByID(s.ctx, req.AccessID)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find intelligent access: %s", err.Error()))
		return
	}

	common.GinSuccess(c, &pb.GetResponse{
		IntelligentAccess: CoverDbAccessToPn(intelligentAccess),
	})
}

// ListHandler finds all intelligent access HTTP handler function
func (s *IntelligentAccessService) ListHandler(c *gin.Context) {
	intelligentAccesses, err := mysql.IntelligentAccessRepo.FindAll(s.ctx)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find all intelligent access: %s", err.Error()))
		return
	}

	var pbIntelligentAccesses []*pb.IntelligentAccess
	for _, intelligentAccess := range intelligentAccesses {
		pbIntelligentAccesses = append(pbIntelligentAccesses, CoverDbAccessToPn(intelligentAccess))
	}

	common.GinSuccess(c, &pb.ListResponse{
		List: pbIntelligentAccesses,
	})
}

func (s *IntelligentAccessService) TestConnectionHandler(c *gin.Context) {
	var req pb.TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}

	db, err := BuildTemporaryPostgresConnection(req.DbHost, int(req.DbPort), req.DbUser, req.DbPassword, req.DbName, true)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to connect to database: %s", err.Error()))
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to ping database: %s", err.Error()))
		return
	}

	common.GinSuccess(c, &pb.TestConnectionResponse{
		Success: true,
		Message: "PostgreSQL database connection successful",
	})
}

func (s *IntelligentAccessService) CheckN8NHandler(c *gin.Context) {
	var req pb.CheckN8NRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}
	if req.AccessID > 0 {
		intelligentAccess, err := mysql.IntelligentAccessRepo.FindByID(s.ctx, req.AccessID)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find intelligent access: %s", err.Error()))
			return
		}
		req.BaseUrl = intelligentAccess.BaseUrl
		req.Username = intelligentAccess.Username
		req.Password = intelligentAccess.Password
	}

	cookie, err := n8n.GetCookieFromLogin(req.BaseUrl, req.Username, req.Password)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to get cookie from login: %s", err.Error()))
		return
	}

	pluginList, err := n8n.GetPluginList(req.BaseUrl, cookie)
	if err != nil {
		common.GinSuccess(c, &pb.CheckN8NResponse{
			LoginStatus:  true,
			PluginStatus: false,
			Message:      fmt.Sprintf("failed to get plugin list: %s", err.Error()),
		})
		return
	}
	for _, plugin := range pluginList.Data {
		if plugin.PackageName == N8NPluginName {
			common.GinSuccess(c, &pb.CheckN8NResponse{
				LoginStatus:  true,
				PluginStatus: true,
				Message:      N8NPluginName + " plugin installed",
			})
			return
		}
	}
	common.GinSuccess(c, &pb.CheckN8NResponse{
		LoginStatus:  true,
		PluginStatus: false,
		Message:      N8NPluginName + " plugin not installed",
	})
}

func (s *IntelligentAccessService) InstallN8NPluginHandler(c *gin.Context) {
	var req pb.InstallN8NPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}
	if req.AccessID > 0 {
		intelligentAccess, err := mysql.IntelligentAccessRepo.FindByID(s.ctx, req.AccessID)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find intelligent access: %s", err.Error()))
			return
		}
		req.BaseUrl = intelligentAccess.BaseUrl
		req.Username = intelligentAccess.Username
		req.Password = intelligentAccess.Password
	}

	cookie, err := n8n.GetCookieFromLogin(req.BaseUrl, req.Username, req.Password)
	if err != nil {
		common.GinSuccess(c, &pb.InstallN8NPluginResponse{
			Success: false,
			Message: "Login error: " + err.Error(),
		})
		return
	}

	_, err = n8n.InstallPackage(req.BaseUrl, N8NPluginName, cookie)
	if err != nil {
		common.GinSuccess(c, &pb.InstallN8NPluginResponse{
			Success: false,
			Message: "Install n8n plugin failed error: " + err.Error(),
		})
		return
	}

	common.GinSuccess(c, &pb.InstallN8NPluginResponse{
		Success: true,
		Message: N8NPluginName + " plugin installed",
	})
}

func BuildTemporaryPostgresConnection(host string, port int, user string, password string, database string, test bool) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, port, user, password, database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %s", err)
	}

	// 合理的连接池配置
	if test {
		db.SetMaxOpenConns(1) // 根据并发需求调整
		db.SetMaxIdleConns(0)
	} else {
		db.SetMaxOpenConns(10) // 根据并发需求调整
		db.SetMaxIdleConns(5)
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}
func (s *IntelligentAccessService) ListUserSpaceHandler(c *gin.Context) {
	var req pb.ListUserSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.GinError(c, i18nresp.CodeBadRequest, err.Error())
		return
	}

	intelligentAccess, err := mysql.IntelligentAccessRepo.FindByID(s.ctx, req.AccessID)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find intelligent access: %s", err.Error()))
		return
	}

	pbUserSpaces := []*pb.UserSpace{}
	if intelligentAccess.AccessType == pb.IntelligentAccessType_Dify.String() || intelligentAccess.AccessType == pb.IntelligentAccessType_DifyEnterprise.String() || intelligentAccess.AccessType == pb.IntelligentAccessType_QAgent.String() {
		// 连接数据库
		sqlDB, err := BuildTemporaryPostgresConnection(intelligentAccess.DbHost, int(intelligentAccess.DbPort), intelligentAccess.DbUser, intelligentAccess.DbPassword, intelligentAccess.DbName, true)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to connect to database: %s", err.Error()))
			return
		}
		defer sqlDB.Close()

		pbUserSpaces, err = GetDifyUserSpace(sqlDB)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to query user spaces: %s", err.Error()))
			return
		}
	}
	if intelligentAccess.AccessType == pb.IntelligentAccessType_COZE.String() {
		pbUserSpaces, err = GetCozeUserSpace(req.Cookie, intelligentAccess)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to query user spaces: %s", err.Error()))
			return
		}
	}
	if intelligentAccess.AccessType == pb.IntelligentAccessType_N8N.String() {
		cookie, err := n8n.GetCookieFromLogin(intelligentAccess.BaseUrl, intelligentAccess.Username, intelligentAccess.Password)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to get n8n cookie: %s", err.Error()))
			return
		}
		pbUserSpaces, err = GetN8NUserSpace(intelligentAccess, cookie)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to query user spaces: %s", err.Error()))
			return
		}
	}

	// 给每个用户回显设置对应实例的 token
	for _, instanceID := range req.InstancesIDs {
		// 获取实例令牌
		tokens, err := mysql.McpTokenRepo.ListByInstanceID(context.Background(), instanceID)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to find mcp instance token: %s", err.Error()))
			return
		}

		if len(tokens) == 0 {
			continue
		}
		// 找到默认 token
		defaultToken := tokens[0]
		for _, token := range tokens {
			usages := []string{}
			_ = json.Unmarshal(token.Usages, &usages)
			for _, usage := range usages {
				if usage == "default" {
					defaultToken = token
					break
				}
			}
		}

		for _, userSpace := range pbUserSpaces {
			userSpace.EncryptPublicKey = ""

			if userSpace.Headers == nil {
				userSpace.Headers = map[string]*pb.HeaderInfo{}
			}

			var findToken *model.McpToken
			if intelligentAccess.AccessType == pb.IntelligentAccessType_Dify.String() || intelligentAccess.AccessType == pb.IntelligentAccessType_DifyEnterprise.String() || intelligentAccess.AccessType == pb.IntelligentAccessType_QAgent.String() {
				findToken = FindToken(tokens, &model.InsertIntelligentInfo{
					SpaceID: userSpace.TenantID,
					UserID:  userSpace.UserID,
				}, req.AccessID)
			} else if intelligentAccess.AccessType == pb.IntelligentAccessType_COZE.String() {
				findToken = FindToken(tokens, &model.InsertIntelligentInfo{
					SpaceID: userSpace.TenantID,
					UserID:  userSpace.UserID,
				}, req.AccessID)
			}

			if findToken != nil {
				headers := map[string]string{}
				_ = json.Unmarshal(findToken.Headers, &headers)

				userSpace.Headers[instanceID] = &pb.HeaderInfo{
					Token:   findToken.Token,
					Headers: headers,
				}
			} else {
				headers := map[string]string{}
				_ = json.Unmarshal(defaultToken.Headers, &headers)
				token, _ := GenerateBearerToken()
				// 没有找到对应 token，生成一个 token
				userSpace.Headers[instanceID] = &pb.HeaderInfo{
					Token:   "Bearer " + token,
					Headers: headers,
				}
			}
		}
	}

	common.GinSuccess(c, &pb.ListUserSpaceResponse{
		UserSpaces: pbUserSpaces,
	})
}

func GenerateBearerToken() (string, error) {
	// 1. 生成 UUID（类似示例中的第一部分）
	uuidObj := uuid.New()
	uuidBase64 := base64.StdEncoding.EncodeToString([]byte(uuidObj.String()))

	// 2. 创建 JSON 负载

	// 3. 序列化 JSON
	payloadJSON, err := json.Marshal(map[string]string{
		"expire_at": fmt.Sprintf("%d", time.Now().Add(24*time.Hour).Unix()*1000), // 毫秒时间戳
		"user_id":   fmt.Sprintf("%d", 1),
		"username":  "admin",
	})
	if err != nil {
		return "", err
	}

	// 4. Base64 编码 JSON
	payloadBase64 := base64.StdEncoding.EncodeToString(payloadJSON)

	// 5. 组合成完整 token
	fullToken := uuidBase64 + payloadBase64

	return fullToken, nil
}

func GetN8NUserSpace(access *model.IntelligentAccess, cookie string) ([]*pb.UserSpace, error) {
	projects, err := n8n.GetProjects(access.BaseUrl, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to get n8n projects: %s", err.Error())
	}
	var pbUserSpaces []*pb.UserSpace
	for _, project := range projects.Data {
		// 判定 Scopes 是否同时存在这三个 credential:create credential:list credential:update 字符串
		if !containsAll(project.Scopes, "credential:create", "credential:list", "credential:update") {
			continue
		}

		pbUserSpaces = append(pbUserSpaces, &pb.UserSpace{
			UserID:           project.ID,
			UserName:         project.Name,
			TenantID:         project.ID,
			TenantName:       project.Name,
			EncryptPublicKey: "",
		})
	}
	return pbUserSpaces, nil
}

func containsAll(arr []string, targets ...string) bool {
	// 创建查找集合
	set := make(map[string]struct{})
	for _, str := range arr {
		set[str] = struct{}{}
	}

	// 检查所有目标字符串是否都在集合中
	for _, target := range targets {
		if _, exists := set[target]; !exists {
			return false
		}
	}
	return true
}

func GetCozeUserSpace(cookie string, access *model.IntelligentAccess) ([]*pb.UserSpace, error) {
	var organizationID = ""
	var err error
	if access.EnterpriseID != "" {
		organizationID, err = coze.GetOrganizationID(cookie, access.EnterpriseID)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization id: %s", err.Error())
		}
	}
	spaceList, err := coze.GetSpaceList(cookie, access.EnterpriseID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get space list: %s", err.Error())
	}

	var pbUserSpaces []*pb.UserSpace
	for _, space := range spaceList {
		if space.SpaceRoleType != 2 && space.SpaceRoleType != 1 {
			continue
		}

		// coze 个人版判定，如果是个人版，只能访问到自己的空间
		if access.SubType == pb.SubType_Person.String() && access.CozeUserID != "" {
			if space.OwnerUserID != access.CozeUserID {
				return nil, fmt.Errorf("failed to get space list, coze user id not match: %s", space.OwnerUserID)
			}
		}

		pbUserSpaces = append(pbUserSpaces, &pb.UserSpace{
			UserID:           space.OwnerUserID,
			TenantID:         space.ID,
			UserName:         space.OwnerUserName,
			TenantName:       space.Name,
			EncryptPublicKey: "",
		})
	}
	return pbUserSpaces, nil
}

func GetDifyUserSpace(sqlDB *sql.DB) ([]*pb.UserSpace, error) {
	userSpaces, err := postgres.GetOwnerTenantAccountJoins(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("failed to query user spaces: %s", err.Error())
	}
	tenants, err := postgres.GetAllTenants(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("failed to query tenants: %s", err.Error())
	}
	accounts, err := postgres.GetAllAccounts(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %s", err.Error())
	}
	tenantMap := make(map[string]postgres.Tenant)
	for _, tenant := range tenants {
		tenantMap[tenant.Id] = tenant
	}
	accountMap := make(map[string]string)
	for _, account := range accounts {
		accountMap[account.Id] = account.Name
	}

	var pbUserSpaces []*pb.UserSpace
	for _, userSpace := range userSpaces {
		pbUserSpaces = append(pbUserSpaces, &pb.UserSpace{
			UserID:           userSpace.AccountId,
			TenantID:         userSpace.TenantId,
			UserName:         accountMap[userSpace.AccountId],
			TenantName:       tenantMap[userSpace.TenantId].Name,
			EncryptPublicKey: tenantMap[userSpace.TenantId].EncryptPublicKey,
		})
	}
	return pbUserSpaces, nil
}

func CoverDbAccessToPn(access *model.IntelligentAccess) *pb.IntelligentAccess {
	return &pb.IntelligentAccess{
		AccessID:     access.ID,
		AccessName:   access.AccessName,
		AccessType:   access.AccessType,
		DbHost:       access.DbHost,
		DbPort:       int32(access.DbPort),
		DbUser:       access.DbUser,
		DbPassword:   access.DbPassword,
		DbName:       access.DbName,
		CreateTime:   access.CreateTime.UnixMilli(),
		UpdateTime:   access.UpdateTime.UnixMilli(),
		SubType:      access.SubType,
		EnterpriseID: access.EnterpriseID,
		BaseUrl:      access.BaseUrl,
		Username:     access.Username,
		Password:     access.Password,
		CozeUserID:   access.CozeUserID,
	}
}
