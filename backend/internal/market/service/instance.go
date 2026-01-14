package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	instancepb "github.com/kymo-mcp/mcpcan/api/market/instance"
	"github.com/kymo-mcp/mcpcan/internal/market/biz"
	"github.com/kymo-mcp/mcpcan/pkg/common"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
	i18nresp "github.com/kymo-mcp/mcpcan/pkg/i18n"
)

// InstanceService struct for instance service
type InstanceService struct {
}

// ContainerInstanceUpdateParams parameters for container instance update
type ContainerInstanceUpdateParams struct {
	Instance                   *model.McpInstance
	ContainerName              string
	ServiceName                string
	ContainerStatus            model.ContainerStatus
	ContainerLastMessage       string
	ContainerCreateOptions     []byte
	ContainerInitTimeoutStopAt int64
	ContainerRunTimeoutStopAt  int64
	Status                     model.InstanceStatus
	ContainerIsReady           bool
}

// NewInstanceService creates a new instance service
func NewInstanceService() *InstanceService {
	return &InstanceService{}
}

// CreateHandler creates instance HTTP handler function
func (s *InstanceService) CreateHandler(c *gin.Context) {
	var req instancepb.CreateRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	result, err := biz.GInstanceBiz.CreateInstance(c.Request.Context(), &req)
	if err != nil {
		// Assuming biz layer returns errors that can be mapped to status codes, or default to 500
		// For now, based on original code, most errors were 500 or 403
		if strings.Contains(err.Error(), "forbidden") {
			common.GinError(c, i18nresp.CodeForbidden, err.Error())
		} else if strings.Contains(err.Error(), "missing required field") {
			common.GinError(c, i18nresp.CodeInternalError, err.Error())
		} else {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to write instance: %s", err.Error()))
		}
		return
	}

	// Return success response
	common.GinSuccess(c, result)
}

// DetailHandler gets instance details HTTP handler function
func (s *InstanceService) DetailHandler(c *gin.Context) {
	var req instancepb.DetailRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Call get instance details handler function
	result, err := s.detail(&req)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to get instance details: %s", err.Error()))
		return
	}

	// Return success response
	common.GinSuccess(c, result)
}

// EditHandler instance edit
func (s *InstanceService) EditHandler(c *gin.Context) {
	var req instancepb.EditRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Get original instance information
	oriInstance, err := biz.GInstanceBiz.GetInstance(req.InstanceId)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to get instance information: %s", err.Error()))
		return
	}
	if oriInstance == nil {
		common.GinError(c, i18nresp.CodeInternalError, "instance does not exist")
		return
	}

	var resp *instancepb.EditResp
	switch oriInstance.AccessType {
	case model.AccessTypeDirect:
		// validate instance name
		if len(req.Name) == 0 {
			common.GinError(c, i18nresp.CodeInternalError, "missing required field: name")
			return
		}
		// validate instance mcpServers
		if len(req.McpServers) == 0 {
			common.GinError(c, i18nresp.CodeInternalError, "missing required field: mcpServers")
			return
		}
		resp, err = biz.GInstanceBiz.UpdateInstanceForDirect(c.Request.Context(), &req, oriInstance)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to edit instance: %s", err.Error()))
			return
		}
	case model.AccessTypeProxy:
		// validate instance name
		if len(req.Name) == 0 {
			common.GinError(c, i18nresp.CodeInternalError, "missing required field: name")
			return
		}
		// validate instance mcpServers
		if len(req.McpServers) == 0 {
			common.GinError(c, i18nresp.CodeInternalError, "missing required field: mcpServers")
			return
		}
		resp, err = biz.GInstanceBiz.UpdateInstanceForProxy(c.Request.Context(), &req, oriInstance)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to edit instance: %s", err.Error()))
			return
		}
	case model.AccessTypeHosting:
		// validate instance name
		if len(req.Name) == 0 {
			common.GinError(c, i18nresp.CodeInternalError, "missing required field: name")
			return
		}
		// validate instance port
		if req.Port <= 0 {
			common.GinError(c, i18nresp.CodeInternalError, "missing required field: port")
			return
		}
		resp, err = biz.GInstanceBiz.UpdateInstanceForHosting(c.Request.Context(), &req, oriInstance)
		if err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to edit instance: %s", err.Error()))
			return
		}
	default:
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("unknown access type: %s", oriInstance.AccessType))
		return
	}

	common.GinSuccess(c, resp)
}

// ListHandler instance list
func (s *InstanceService) ListHandler(c *gin.Context) {
	var req instancepb.ListRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Use InstanceService to handle request
	result, err := s.list(&req)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to get instance list: %s", err.Error()))
		return
	}

	common.GinSuccess(c, result)
}

// DisabledHandler disable instance handler
func (s *InstanceService) DisabledHandler(c *gin.Context) {
	var req instancepb.DisabledRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Use InstanceService to handle request
	result, err := s.disable(c.Request.Context(), &req)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, err.Error())
		return
	}

	common.GinSuccess(c, result)
}

// RestartHandler restart instance handler
func (s *InstanceService) RestartHandler(c *gin.Context) {
	var req instancepb.RestartRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}
	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Use InstanceService to handle request
	result, err := biz.GInstanceBiz.RestartInstance(c.Request.Context(), &req)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, err.Error())
		return
	}

	common.GinSuccess(c, result)
}

// DeleteHandler delete instance handler
func (s *InstanceService) DeleteHandler(c *gin.Context) {
	var req instancepb.DeleteRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Use InstanceService to handle request
	result, err := biz.GInstanceBiz.DeleteInstance(req.InstanceId)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, err.Error())
		return
	}

	common.GinSuccess(c, result)
}

// StatusHandler query instance status handler
func (s *InstanceService) StatusHandler(c *gin.Context) {
	var req instancepb.GetStatusRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Use InstanceService to handle request
	result, err := biz.GInstanceBiz.GetStatus(c.Request.Context(), &req)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, err.Error())
		return
	}

	common.GinSuccess(c, result)
}

// LogsHandler get managed instance logs handler
func (s *InstanceService) LogsHandler(c *gin.Context) {
	var req instancepb.LogsRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Use InstanceService to handle request
	result, err := biz.GInstanceBiz.GetLogs(c.Request.Context(), &req)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, err.Error())
		return
	}

	common.GinSuccess(c, result)
}

// create writes instance method
// func (s *InstanceService) create(req *instancepb.CreateRequest) (*instancepb.CreateResp, error) {
//
// 	// Generate instance ID (UUID)
// 	instanceID := uuid.New().String()
//
// 	// Hosting mode, Stdio protocol
// 	switch req.AccessType {
// 	case instancepb.AccessType_DIRECT:
// 		return s.createInstanceDirectMode(req, instanceID)
// 	case instancepb.AccessType_PROXY:
// 		return s.createInstanceProxyMode(req, instanceID)
// 	case instancepb.AccessType_HOSTING:
// 		return s.createInstanceHosting(req, instanceID)
// 	default:
// 		return nil, fmt.Errorf("unsupported access type: %v", req.AccessType)
// 	}
// }

// detail gets instance details
func (s *InstanceService) detail(req *instancepb.DetailRequest) (*instancepb.DetailResp, error) {
	// Get instance information
	instance, err := biz.GInstanceBiz.GetInstance(req.InstanceId)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance information: %v", err)
	}
	if instance == nil {
		return nil, fmt.Errorf("instance does not exist")
	}

	// Convert access type
	pbAccessType, err := common.ConvertToProtoAccessType(instance.AccessType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access type: %w", err)
	}

	// Convert MCP protocol type
	pbMcpProtocol, err := common.ConvertToProtoMcpProtocol(instance.McpProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to convert MCP protocol type: %w", err)
	}

	// Build response
	resp := &instancepb.DetailResp{
		InstanceId:     instance.InstanceID,
		Name:           instance.InstanceName,
		Status:         string(instance.Status),
		AccessType:     pbAccessType,
		McpProtocol:    pbMcpProtocol,
		Notes:          instance.Notes,
		IconPath:       instance.IconPath,
		EnabledToken:   instance.EnabledToken,
		OpenapiBaseUrl: instance.OpenapiBaseUrl,
	}

	// Add specific fields based on access type
	switch instance.AccessType {
	case model.AccessTypeHosting:
		resp.PackageId = instance.PackageID
		resp.EnvironmentId = int32(instance.EnvironmentID)
		resp.McpServerId = instance.McpServerID
		resp.TemplateId = int32(instance.TemplateID)
		resp.ImgAddress = instance.ImgAddr
		resp.McpServers = string(instance.SourceConfig)
		resp.Port = instance.Port
		resp.InitScript = instance.InitScript
		resp.Command = instance.Command
		resp.ServicePath = instance.ServicePath
		resp.ContainerName = instance.ContainerName
		resp.ContainerServiceName = instance.ContainerServiceName
		resp.ContainerStatus = string(instance.ContainerStatus)
		resp.ContainerLastMessage = instance.ContainerLastMessage
		resp.ContainerIsReady = instance.ContainerIsReady

		// Convert environment variables
		if len(instance.EnvironmentVariables) > 0 {
			envVarsMap := make(map[string]string)
			if err := json.Unmarshal(instance.EnvironmentVariables, &envVarsMap); err == nil {
				resp.EnvironmentVariables = envVarsMap
			}
		}

		// Convert volume mounts
		if len(instance.VolumeMounts) > 0 {
			var volumeMounts []*instancepb.VolumeMount
			if err := json.Unmarshal(instance.VolumeMounts, &volumeMounts); err == nil {
				resp.VolumeMounts = volumeMounts
			}
		}
	case model.AccessTypeDirect, model.AccessTypeProxy:
		// For direct and proxy mode, add MCP servers configuration
		if len(instance.SourceConfig) > 0 {
			resp.McpServers = string(instance.SourceConfig)
		}
	}

	return resp, nil
}

func (s *InstanceService) list(req *instancepb.ListRequest) (*instancepb.ListResp, error) {
	// Parameter validation
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = int32(common.DefaultPageSize)
	}
	if pageSize > int32(common.MaxPageSize) {
		pageSize = int32(common.MaxPageSize)
	}

	// Build filter conditions
	filters := make(map[string]interface{})
	if req.InstanceName != "" {
		filters["instanceName"] = req.InstanceName
	}
	if req.EnvironmentId > 0 {
		filters["environmentId"] = req.EnvironmentId
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.ContainerStatus != "" {
		filters["containerStatus"] = req.ContainerStatus
	}
	if req.AccessType > 0 {
		accessType, err := common.ConvertToModelAccessType(req.AccessType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert access type: %w", err)
		}
		filters["accessType"] = accessType
	}
	if req.McpProtocol > 0 {
		mcpProtocol, err := common.ConvertToModelMcpProtocol(req.McpProtocol)
		if err != nil {
			return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
		}
		filters["mcpProtocol"] = mcpProtocol
	}

	// Sort parameters
	sortBy := "createdAt"
	sortOrder := "desc"

	return biz.GInstanceBiz.ListInstance(page, pageSize, filters, sortBy, sortOrder)
}

// disable disables an instance
func (s *InstanceService) disable(ctx context.Context, req *instancepb.DisabledRequest) (*instancepb.DisabledResp, error) {
	// Disable the instance and set deletion time
	msg, err := biz.GInstanceBiz.DisableInstance(ctx, req.InstanceId)
	if err != nil {
		return nil, fmt.Errorf("failed to disable instance: %v", err)
	}

	return &instancepb.DisabledResp{Message: msg}, nil
}

// createInstanceDirectMode direct connection mode handler function
// func (s *InstanceService) createInstanceDirectMode(req *instancepb.CreateRequest, instanceID string) (*instancepb.CreateResp, error) {
// 	accessType, err := common.ConvertToModelAccessType(req.AccessType)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert access type: %w", err)
// 	}
// 	mcpProtocol, err := common.ConvertToModelMcpProtocol(req.McpProtocol)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
// 	}
// 	sourceType, err := common.ConvertToModelSourceType(req.SourceType)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert source type: %w", err)
// 	}
// 	if len(req.McpServers) == 0 {
// 		return nil, fmt.Errorf("missing required field: mcpServers")
// 	}
// 	// Validate MCP configuration format
// 	validationResult, err := utils.ValidateMcpConfig([]byte(req.McpServers))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to validate mcp servers: %w", err)
// 	}
// 	if !validationResult.IsValid {
// 		return nil, fmt.Errorf("mcp servers config is invalid: %s", validationResult.ErrorMessage)
// 	}
// 	if validationResult.Url == "" {
// 		return nil, fmt.Errorf("mcp servers config is invalid: url is empty")
// 	}
// 	if validationResult.ProtocolType != string(mcpProtocol) {
// 		return nil, fmt.Errorf("mcp servers config is invalid: protocol type is %s, expected %s", validationResult.ProtocolType, mcpProtocol)
// 	}
//
// 	sourceConfig := json.RawMessage([]byte(req.McpServers))
// 	// Create new instance record
// 	instance := &model.McpInstance{
// 		InstanceID:   instanceID,
// 		InstanceName: req.Name,
// 		AccessType:   accessType,
// 		McpProtocol:  mcpProtocol,
// 		SourceType:   sourceType,
// 		SourceConfig: sourceConfig,
// 		Status:       model.InstanceStatusActive,
// 		IconPath:     req.IconPath,         // Add iconPath field handling
// 		Notes:        req.Notes,            // Add notes field handling
// 		McpServerID:  req.McpServerId,      // Add mcpServerId field handling
// 		TemplateID:   uint(req.TemplateId), // Add templateId field handling
// 	}
//
// 	// Save instance to database
// 	if err := biz.GInstanceBiz.CreateInstance(instance); err != nil {
// 		return nil, fmt.Errorf("failed to create instance: %w", err)
// 	}
//
// 	return &instancepb.CreateResp{
// 		InstanceId:  instanceID,
// 		Name:        req.Name,
// 		Status:      string(model.InstanceStatusActive),
// 		AccessType:  req.AccessType,
// 		McpProtocol: req.McpProtocol,
// 	}, nil
// }

// createInstanceProxyMode proxy mode handler function
// func (s *InstanceService) createInstanceProxyMode(req *instancepb.CreateRequest, instanceID string) (*instancepb.CreateResp, error) {
// 	accessType, err := common.ConvertToModelAccessType(req.AccessType)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert access type: %w", err)
// 	}
// 	mcpProtocol, err := common.ConvertToModelMcpProtocol(req.McpProtocol)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
// 	}
// 	sourceType, err := common.ConvertToModelSourceType(req.SourceType)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert source type: %w", err)
// 	}
// 	if len(req.McpServers) == 0 {
// 		return nil, fmt.Errorf("missing required field: mcpServers")
// 	}
// 	// Validate MCP configuration format
// 	validationResult, err := utils.ValidateMcpConfig([]byte(req.McpServers))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to validate mcp servers: %w", err)
// 	}
// 	if !validationResult.IsValid {
// 		return nil, fmt.Errorf("mcp servers config is invalid: %s", validationResult.ErrorMessage)
// 	}
// 	if validationResult.Url == "" {
// 		return nil, fmt.Errorf("mcp servers config is invalid: url is empty")
// 	}
// 	if validationResult.ProtocolType != string(mcpProtocol) {
// 		return nil, fmt.Errorf("mcp servers config is invalid: protocol type is %s, expected %s", validationResult.ProtocolType, mcpProtocol)
// 	}
// 	proxyProtocol := mcpProtocol
// 	publicProxyPath := biz.GInstanceBiz.CreatePublicProxyPath(instanceID, proxyProtocol)
//
// 	// Create new instance record
// 	instance := &model.McpInstance{
// 		InstanceID:      instanceID,
// 		InstanceName:    req.Name,
// 		AccessType:      accessType,
// 		McpProtocol:     mcpProtocol,
// 		SourceType:      sourceType,
// 		SourceConfig:    json.RawMessage([]byte(req.McpServers)),
// 		Status:          model.InstanceStatusActive,
// 		IconPath:        req.IconPath,         // Add iconPath field handling
// 		Notes:           req.Notes,            // Add notes field handling
// 		McpServerID:     req.McpServerId,      // Add mcpServerId field handling
// 		TemplateID:      uint(req.TemplateId), // Add templateId field handling
// 		EnabledToken:    req.EnabledToken,
// 		PublicProxyPath: publicProxyPath,
// 		ProxyProtocol:   proxyProtocol,
// 	}
//
// 	if len(req.Tokens) > 0 {
// 		// add instance id to tokens
// 		for _, token := range req.Tokens {
// 			token.InstanceId = instanceID
// 		}
// 		if err := biz.GInstanceBiz.SaveTokensForInstance(c.Request.Context(), req.Tokens); err != nil {
// 			return nil, fmt.Errorf("failed to save tokens for instance: %w", err)
// 		}
// 	}
//
// 	// Save instance to database
// 	if err := biz.GInstanceBiz.CreateInstance(instance); err != nil {
// 		return nil, fmt.Errorf("failed to create instance: %w", err)
// 	}
//
// 	return &instancepb.CreateResp{
// 		InstanceId:  instanceID,
// 		Name:        req.Name,
// 		Status:      string(model.InstanceStatusActive),
// 		AccessType:  req.AccessType,
// 		McpProtocol: req.McpProtocol,
// 	}, nil
// }

// createInstanceHosting Hosting mode handler function
// func (s *InstanceService) createInstanceHosting(req *instancepb.CreateRequest, instanceID string) (*instancepb.CreateResp, error) {
// 	// Validate timeout parameters
// 	if err := s.validateTimeoutParams(int(req.StartupTimeout), int(req.RunningTimeout)); err != nil {
// 		return nil, fmt.Errorf("parameter validation failed: %w", err)
// 	}
// 	mcpProtocol, err := common.ConvertToModelMcpProtocol(req.McpProtocol)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
// 	}
// 	sourceType, err := common.ConvertToModelSourceType(req.SourceType)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert source type: %w", err)
// 	}
//
// 	if req.Port <= 0 {
// 		return nil, fmt.Errorf("missing required field: port")
// 	}
// 	// Validate environment ID
// 	if req.EnvironmentId == 0 {
// 		return nil, fmt.Errorf("hosting type instance requires environment ID")
// 	}
// 	if req.ImgAddress == "" {
// 		return nil, fmt.Errorf("missing required field: imgAddress")
// 	}
// 	if mcpProtocol == model.McpProtocolStdio {
// 		mcpServers := req.McpServers
// 		if len(mcpServers) == 0 {
// 			return nil, fmt.Errorf("mcp servers config is empty")
// 		}
// 		reqMcpResult, err2 := utils.ValidateMcpConfig([]byte(mcpServers))
// 		if err2 != nil {
// 			return nil, fmt.Errorf("failed to validate mcp servers: %w", err2)
// 		}
// 		if !reqMcpResult.IsValid {
// 			return nil, fmt.Errorf("mcp servers config is invalid: %s", reqMcpResult.ErrorMessage)
// 		}
// 		if !reqMcpResult.HasCommand {
// 			return nil, fmt.Errorf("mcp servers config is invalid: command is required")
// 		}
// 	}
// 	containerOptions, err := biz.GContainerBiz.BuildContainerOptions(c.Request.Context(), instanceID, mcpProtocol, req.McpServers, req.PackageId, req.Port,
// 		req.InitScript, req.Command, req.ImgAddress, req.EnvironmentVariables, req.VolumeMounts, int32(req.StartupTimeout), int32(req.RunningTimeout))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to build container options: %w", err)
// 	}
// 	err = biz.GContainerBiz.CreateContainer(containerOptions, req.EnvironmentId, req.StartupTimeout)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create container: %w", err)
// 	}
//
// 	// Create new instance record
// 	containerCreateOptions, err := common.MarshalAndAssignConfig(containerOptions)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal container create options: %w", err)
// 	}
// 	evs, err := common.MarshalAndAssignConfig(req.EnvironmentVariables)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal environment variables: %w", err)
// 	}
// 	vms, err := common.MarshalAndAssignConfig(req.VolumeMounts)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal volume mounts: %w", err)
// 	}
// 	// Create target configuration
// 	proxyProtocol := mcpProtocol
// 	publicProxyPath := ""
// 	containerServiceURL := ""
// 	switch mcpProtocol {
// 	case model.McpProtocolStdio:
// 		proxyProtocol = model.McpProtocolStreamableHttp
// 		publicProxyPath = biz.GInstanceBiz.CreatePublicProxyPath(instanceID, proxyProtocol)
// 		containerServiceURL = fmt.Sprintf("http://%s:%d/%s", containerOptions.ServiceName, containerOptions.Port, "mcp")
// 	case model.McpProtocolSSE, model.McpProtocolStreamableHttp:
// 		proxyProtocol = mcpProtocol
// 		publicProxyPath = biz.GInstanceBiz.CreatePublicProxyPath(instanceID, proxyProtocol)
// 		containerServiceURL = fmt.Sprintf("http://%s:%d/%s", containerOptions.ServiceName, containerOptions.Port, strings.Trim(req.ServicePath, "/"))
// 	default:
// 		return nil, fmt.Errorf("unsupported mcp protocol: %v", mcpProtocol)
// 	}
// 	instance := &model.McpInstance{
// 		InstanceID:             instanceID,
// 		InstanceName:           req.Name,
// 		AccessType:             model.AccessTypeHosting,
// 		McpProtocol:            mcpProtocol,
// 		Status:                 model.InstanceStatusActive,
// 		PackageID:              req.PackageId,
// 		ContainerStatus:        model.ContainerStatusPending,
// 		EnvironmentID:          uint(req.EnvironmentId),
// 		SourceType:             sourceType,
// 		McpServerID:            req.McpServerId,
// 		TemplateID:             uint(req.TemplateId),
// 		EnabledToken:           req.EnabledToken,
// 		ImgAddr:                req.ImgAddress,
// 		Port:                   req.Port,
// 		InitScript:             req.InitScript,
// 		Command:                req.Command,
// 		EnvironmentVariables:   evs,
// 		VolumeMounts:           vms,
// 		ContainerName:          containerOptions.ContainerName,
// 		ContainerServiceName:   containerOptions.ServiceName,
// 		ContainerIsReady:       false,
// 		ContainerCreateOptions: containerCreateOptions,
// 		ContainerLastMessage:   "container is pending",
// 		ContainerServiceURL:    containerServiceURL,
// 		StartupTimeout:         int64(req.StartupTimeout),
// 		RunningTimeout:         int64(req.RunningTimeout),
// 		SourceConfig:           json.RawMessage(req.McpServers),
// 		ServicePath:            req.ServicePath,
// 		Notes:                  req.Notes,
// 		IconPath:               req.IconPath,
// 		PublicProxyPath:        publicProxyPath,
// 		ProxyProtocol:          proxyProtocol,
// 	}
// 	if len(req.Tokens) > 0 {
// 		// add instance id to tokens
// 		for _, token := range req.Tokens {
// 			token.InstanceId = instanceID
// 		}
// 		if err := biz.GInstanceBiz.SaveTokensForInstance(c.Request.Context(), req.Tokens); err != nil {
// 			return nil, fmt.Errorf("failed to save tokens for instance: %w", err)
// 		}
// 	}
// 	// Save instance to database
// 	if err := biz.GInstanceBiz.CreateInstance(instance); err != nil {
// 		return nil, fmt.Errorf("failed to create instance: %w", err)
// 	}
//
// 	return &instancepb.CreateResp{
// 		InstanceId:  instanceID,
// 		Name:        req.Name,
// 		Status:      string(model.InstanceStatusActive),
// 		AccessType:  req.AccessType,
// 		McpProtocol: req.McpProtocol,
// 	}, nil
// }

// validateTimeoutParams validates timeout parameters
// func (s *InstanceService) validateTimeoutParams(startupTimeout, runningTimeout int) error {
// 	// Startup timeout validation
// 	if startupTimeout < 0 {
// 		return fmt.Errorf("startup timeout cannot be negative")
// 	}
// 	if startupTimeout > 0 && startupTimeout < 30 {
// 		return fmt.Errorf("startup timeout cannot be less than 30 seconds")
// 	}
// 	if startupTimeout > 3600 {
// 		return fmt.Errorf("startup timeout cannot exceed 3600 seconds")
// 	}
//
// 	// Running timeout validation
// 	if runningTimeout < 0 {
// 		return fmt.Errorf("running timeout cannot be negative")
// 	}
// 	if runningTimeout > 0 && runningTimeout < 60 {
// 		return fmt.Errorf("running timeout cannot be less than 60 seconds")
// 	}
// 	if runningTimeout > 86400 {
// 		return fmt.Errorf("running timeout cannot exceed 86400 seconds")
// 	}
//
// 	return nil
// }

// TokenControlHandler controls token enable/disable status
func (s *InstanceService) TokenControlHandler(c *gin.Context) {
	var req instancepb.TokenControlRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Get original instance information
	instance, err := biz.GInstanceBiz.GetInstance(req.InstanceId)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to get instance information: %s", err.Error()))
		return
	}
	if instance == nil {
		common.GinError(c, i18nresp.CodeInternalError, "instance does not exist")
		return
	}

	// Update enabledToken field
	instance.EnabledToken = req.EnabledToken
	if err := mysql.McpInstanceRepo.Update(c.Request.Context(), instance); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to update instance token status: %s", err.Error()))
		return
	}

	// Return success response
	common.GinSuccess(c, &instancepb.TokenControlResponse{
		InstanceId:   instance.InstanceID,
		EnabledToken: instance.EnabledToken,
		Message:      "Token status updated successfully",
	})
}

// TokenEditHandler edits tokens with validation and overwrite functionality
func (s *InstanceService) TokenEditHandler(c *gin.Context) {
	var req instancepb.TokenEditRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Delegate to biz for upsert logic
	if err := biz.GInstanceBiz.SaveTokensForInstance(c.Request.Context(), req.Tokens); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to update instance tokens: %s", err.Error()))
		return
	}
	common.GinSuccess(c, &instancepb.TokenEditResponse{
		Message: "Tokens updated successfully",
	})
}

// TokenListByInstanceIDHandler lists tokens for instance
func (s *InstanceService) TokenListByInstanceIDHandler(c *gin.Context) {
	var req instancepb.TokenListByInstanceIDRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}
	resp, err := biz.GInstanceBiz.TokenListByInstanceID(&req)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to list tokens: %s", err.Error()))
		return
	}
	common.GinSuccess(c, resp)
}

// TokenDeleteHandler delete a token
func (s *InstanceService) TokenDeleteHandler(c *gin.Context) {
	var req instancepb.TokenDeleteRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}
	if req.Id == 0 {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: id")
		return
	}

	if err := biz.GInstanceBiz.DeleteTokenByID(c, req.Id); err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to delete token: %s", err.Error()))
		return
	}
	resp := &instancepb.TokenDeleteResponse{Message: "Token deleted successfully"}
	common.GinSuccess(c, resp)
}

// CreateOpenapiHandler creates openapi instance HTTP handler function
func (s *InstanceService) CreateOpenapiHandler(c *gin.Context) {
	var req instancepb.CreateOpenapiRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	err := s.checkSaveOpenapiInstanceRequest(req.Name, req.OpenapiFileID, req.ChooseOpenapiFileID, req.Tokens)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to write instance: %s", err.Error()))
		return
	}

	result, err := biz.GInstanceBiz.CreateOpenapiInstance(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			common.GinError(c, i18nresp.CodeForbidden, err.Error())
		} else {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to write instance: %s", err.Error()))
		}
		return
	}

	// Return success response
	common.GinSuccess(c, result)
}

// UpdateOpenapiHandler update openapi instance HTTP handler function
func (s *InstanceService) UpdateOpenapiHandler(c *gin.Context) {
	var req instancepb.UpdateOpenapiRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	// Validate required fields
	if req.InstanceId == "" {
		common.GinError(c, i18nresp.CodeInternalError, "missing required field: instanceId")
		return
	}

	// Get original instance information
	oriInstance, err := biz.GInstanceBiz.GetInstance(req.InstanceId)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to get instance information: %s", err.Error()))
		return
	}
	if oriInstance == nil {
		common.GinError(c, i18nresp.CodeInternalError, "instance does not exist")
		return
	}

	err = s.checkSaveOpenapiInstanceRequest(req.Name, req.OpenapiFileID, req.ChooseOpenapiFileID, nil)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to write instance: %s", err.Error()))
		return
	}

	resp, err := biz.GInstanceBiz.UpdateInstanceForOpenapi(c.Request.Context(), &req, oriInstance)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to edit instance: %s", err.Error()))
		return
	}

	common.GinSuccess(c, resp)
}

func (s *InstanceService) checkSaveOpenapiInstanceRequest(name string, openapiFileID string, chooseOpenapiFileID string, tokens []*instancepb.McpToken) error {
	if name == "" {
		return fmt.Errorf("missing required field: name")
	}
	// Validate environment ID
	if openapiFileID == "" {
		return fmt.Errorf("missing required field: openapiFileID")
	}
	if chooseOpenapiFileID == "" {
		return fmt.Errorf("openapi type instance requires choose openapi file ID")
	}
	return nil
}

func (s *InstanceService) ListToolsHandler(c *gin.Context) {
	var req instancepb.ListToolsRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}

	tools, err := biz.GInstanceBiz.ListTools(c.Request.Context(), req.InstanceId, req.Domain)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to list tools: %s", err.Error()))
		return
	}

	common.GinSuccess(c, tools)
}

func (s *InstanceService) CallToolHandler(c *gin.Context) {
	var req instancepb.CallToolRequest
	if err := common.BindAndValidate(c, &req); err != nil {
		return
	}
	var args any
	if req.Arguments != "" {
		if err := json.Unmarshal([]byte(req.Arguments), &args); err != nil {
			common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to unmarshal arguments: %s", err.Error()))
			return
		}
	}

	resp, err := biz.GInstanceBiz.CallTool(c.Request.Context(), req.InstanceId, req.ToolName, args, req.Domain)
	if err != nil {
		common.GinError(c, i18nresp.CodeInternalError, fmt.Sprintf("failed to call tool: %s", err.Error()))
		return
	}

	common.GinSuccess(c, resp)
}
