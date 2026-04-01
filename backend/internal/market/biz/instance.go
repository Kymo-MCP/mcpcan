package biz

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kymo-mcp/mcpcan/internal/market/config"
	"github.com/kymo-mcp/mcpcan/pkg/common"
	"github.com/kymo-mcp/mcpcan/pkg/container"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
	"github.com/kymo-mcp/mcpcan/pkg/redis"
	"github.com/kymo-mcp/mcpcan/pkg/utils"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"gorm.io/gorm"

	instancepb "github.com/kymo-mcp/mcpcan/api/market/instance"
)

// InstanceBiz instance data processing layer
type InstanceBiz struct {
	ctx context.Context
}

// GInstanceBiz global instance data processing layer instance
var GInstanceBiz *InstanceBiz

func init() {
	GInstanceBiz = NewInstanceBiz(context.Background())
}

// NewInstanceBiz create instance data processing layer instance
func NewInstanceBiz(ctx context.Context) *InstanceBiz {
	return &InstanceBiz{
		ctx: ctx,
	}
}

// CreateInstance creates a new instance
func (biz *InstanceBiz) CreateInstance(ctx context.Context, req *instancepb.CreateRequest) (*instancepb.CreateResp, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("missing required field: name")
	}
	// Demo mode guard: enforce instance limit
	if config.IsDemoMode() {
		instances, err := mysql.McpInstanceRepo.FindByStatus(ctx, model.InstanceStatusActive)
		if err != nil {
			return nil, fmt.Errorf("failed to count active instances: %s", err.Error())
		}
		if len(instances) >= config.GetDemoMaxInstances() {
			return nil, fmt.Errorf("operation forbidden in demo mode: instance limit reached, max: %d", config.GetDemoMaxInstances())
		}
	}

	// Generate instance ID (UUID)
	instanceID := uuid.New().String()

	// Hosting mode, Stdio protocol
	switch req.AccessType {
	case instancepb.AccessType_DIRECT:
		return biz.createInstanceDirectMode(ctx, req, instanceID)
	case instancepb.AccessType_PROXY:
		return biz.createInstanceProxyMode(ctx, req, instanceID)
	case instancepb.AccessType_HOSTING:
		return biz.createInstanceHosting(ctx, req, instanceID)
	default:
		return nil, fmt.Errorf("unsupported access type: %v", req.AccessType)
	}
}

// createInstanceDirectMode direct connection mode handler function
func (biz *InstanceBiz) createInstanceDirectMode(ctx context.Context, req *instancepb.CreateRequest, instanceID string) (*instancepb.CreateResp, error) {
	accessType, err := common.ConvertToModelAccessType(req.AccessType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access type: %w", err)
	}
	mcpProtocol, err := common.ConvertToModelMcpProtocol(req.McpProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
	}
	sourceType, err := common.ConvertToModelSourceType(req.SourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert source type: %w", err)
	}
	if len(req.McpServers) == 0 {
		return nil, fmt.Errorf("missing required field: mcpServers")
	}
	// Validate MCP configuration format
	validationResult, err := utils.ValidateMcpConfig([]byte(req.McpServers))
	if err != nil {
		return nil, fmt.Errorf("failed to validate mcp servers: %w", err)
	}
	if !validationResult.IsValid {
		return nil, fmt.Errorf("mcp servers config is invalid: %s", validationResult.ErrorMessage)
	}
	if validationResult.Url == "" {
		return nil, fmt.Errorf("mcp servers config is invalid: url is empty")
	}
	if validationResult.ProtocolType != string(mcpProtocol) {
		return nil, fmt.Errorf("mcp servers config is invalid: protocol type is %s, expected %s", validationResult.ProtocolType, mcpProtocol)
	}

	sourceConfig := json.RawMessage([]byte(req.McpServers))
	// Create new instance record
	instance := &model.McpInstance{
		InstanceID:    instanceID,
		InstanceName:  req.Name,
		AccessType:    accessType,
		McpProtocol:   mcpProtocol,
		SourceType:    sourceType,
		SourceConfig:  sourceConfig,
		Status:        model.InstanceStatusActive,
		IconPath:      req.IconPath,         // Add iconPath field handling
		Notes:         req.Notes,            // Add notes field handling
		McpServerID:   req.McpServerId,      // Add mcpServerId field handling
		TemplateID:    uint(req.TemplateId), // Add templateId field handling
		EnvironmentID: uint(req.EnvironmentId),
		EnabledToken:  req.EnabledToken,
	}
	if req.Headers != nil && len(req.Headers) > 0 {
		if b, err := json.Marshal(req.Headers); err == nil {
			instance.Headers = b
		}
	}

	tokensToSave := biz.buildCreateTokens(instanceID, req.EnabledToken, req.Tokens)
	if len(tokensToSave) > 0 {
		if err := biz.SaveTokensForInstance(ctx, tokensToSave); err != nil {
			return nil, fmt.Errorf("failed to save tokens for instance: %w", err)
		}
	}

	// Save instance to database
	if err := biz.CreateInstanceRecord(ctx, instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return &instancepb.CreateResp{
		InstanceId:  instanceID,
		Name:        req.Name,
		Status:      string(model.InstanceStatusActive),
		AccessType:  req.AccessType,
		McpProtocol: req.McpProtocol,
	}, nil
}

// CreateOpenapiInstance creates a new openapi instance
func (biz *InstanceBiz) CreateOpenapiInstance(ctx context.Context, req *instancepb.CreateOpenapiRequest) (*instancepb.CreateResp, error) {
	// Demo mode guard: enforce instance limit
	if config.IsDemoMode() {
		instances, err := mysql.McpInstanceRepo.FindByStatus(ctx, model.InstanceStatusActive)
		if err != nil {
			return nil, fmt.Errorf("failed to count active instances: %s", err.Error())
		}
		if len(instances) >= config.GetDemoMaxInstances() {
			return nil, fmt.Errorf("operation forbidden in demo mode: instance limit reached")
		}
	}

	_, err := mysql.McpOpenapiPackageRepo.FindByOpenapiFileID(ctx, req.OpenapiFileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get openapi file information: %s", err)
	}
	chooseOpenapiFileInfo, err := mysql.McpOpenapiPackageRepo.FindByOpenapiFileID(ctx, req.ChooseOpenapiFileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get choose openapi file information: %s", err)
	}
	if chooseOpenapiFileInfo.BaseOpenapiFileID != req.OpenapiFileID {
		return nil, fmt.Errorf("failed to get openapi file information")
	}

	instanceID := uuid.New().String()
	containerOptions, err := GContainerBiz.BuildOpenapiContainerOptions(ctx, instanceID, chooseOpenapiFileInfo.OpenapiFileID, common.GetMcpHostingPort(), 0, 0, req.OpenapiBaseUrl, req.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to build container options: %w", err)
	}
	err = GContainerBiz.CreateContainer(ctx, containerOptions, req.EnvironmentId, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}
	// Create new instance record
	containerCreateOptions, err := common.MarshalAndAssignConfig(containerOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal container create options: %w", err)
	}

	instance := &model.McpInstance{
		InstanceID:             instanceID,
		InstanceName:           req.Name,
		AccessType:             model.AccessTypeHosting,
		McpProtocol:            model.McpProtocolStreamableHttp,
		Status:                 model.InstanceStatusActive,
		PackageID:              chooseOpenapiFileInfo.OpenapiFileID,
		ContainerStatus:        model.ContainerStatusPending,
		EnvironmentID:          uint(req.EnvironmentId),
		SourceType:             model.SourceTypeOpenapi,
		McpServerID:            "",
		TemplateID:             0,
		EnabledToken:           req.EnabledToken,
		ImgAddr:                containerOptions.ImageName,
		Port:                   common.GetMcpHostingPort(),
		InitScript:             "",
		Command:                containerOptions.CommandArgs[0],
		EnvironmentVariables:   nil,
		VolumeMounts:           nil,
		ContainerName:          containerOptions.ContainerName,
		ContainerServiceName:   containerOptions.ServiceName,
		ContainerIsReady:       false,
		ContainerCreateOptions: containerCreateOptions,
		ContainerLastMessage:   "container is pending",
		ContainerServiceURL:    biz.getContainerServiceURL(containerOptions.ContainerName, containerOptions.ServiceName),
		StartupTimeout:         int64(0),
		RunningTimeout:         int64(0),
		SourceConfig:           nil,
		ServicePath:            "/mcp",
		Notes:                  req.Notes,
		IconPath:               req.IconPath,
		PublicProxyPath:        biz.CreatePublicProxyPath(instanceID, model.McpProtocolStreamableHttp),
		ProxyProtocol:          model.McpProtocolStreamableHttp,
		OpenapiBaseUrl:         req.OpenapiBaseUrl,
	}
	if req.Headers != nil && len(req.Headers) > 0 {
		if b, err := json.Marshal(req.Headers); err == nil {
			instance.Headers = b
		}
	}

	tokensToSave := biz.buildCreateTokens(instanceID, req.EnabledToken, req.Tokens)
	if len(tokensToSave) > 0 {
		if err := biz.SaveTokensForInstance(ctx, tokensToSave); err != nil {
			return nil, fmt.Errorf("failed to save tokens for instance: %w", err)
		}
	}

	// Save instance to database
	if err := biz.CreateInstanceRecord(ctx, instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	// Return success response
	return &instancepb.CreateResp{
		InstanceId:  instanceID,
		Name:        req.Name,
		Status:      string(model.InstanceStatusActive),
		AccessType:  instancepb.AccessType_HOSTING,
		McpProtocol: instancepb.McpProtocol_STREAMABLE_HTTP,
	}, nil
}

// createInstanceProxyMode proxy mode handler function
func (biz *InstanceBiz) createInstanceProxyMode(ctx context.Context, req *instancepb.CreateRequest, instanceID string) (*instancepb.CreateResp, error) {
	accessType, err := common.ConvertToModelAccessType(req.AccessType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access type: %w", err)
	}
	mcpProtocol, err := common.ConvertToModelMcpProtocol(req.McpProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
	}
	sourceType, err := common.ConvertToModelSourceType(req.SourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert source type: %w", err)
	}
	if len(req.McpServers) == 0 {
		return nil, fmt.Errorf("missing required field: mcpServers")
	}
	// Validate MCP configuration format
	validationResult, err := utils.ValidateMcpConfig([]byte(req.McpServers))
	if err != nil {
		return nil, fmt.Errorf("failed to validate mcp servers: %w", err)
	}
	if !validationResult.IsValid {
		return nil, fmt.Errorf("mcp servers config is invalid: %s", validationResult.ErrorMessage)
	}
	if validationResult.Url == "" {
		return nil, fmt.Errorf("mcp servers config is invalid: url is empty")
	}
	if validationResult.ProtocolType != string(mcpProtocol) {
		return nil, fmt.Errorf("mcp servers config is invalid: protocol type is %s, expected %s", validationResult.ProtocolType, mcpProtocol)
	}
	proxyProtocol := mcpProtocol
	publicProxyPath := biz.CreatePublicProxyPath(instanceID, proxyProtocol)
	environmentID, err := biz.resolveProxyEnvironmentID(ctx, req.EnvironmentId)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve proxy runtime environment: %w", err)
	}

	// DEBUG: 临时日志，验证 EnvironmentId 接收值
	fmt.Printf("[DEBUG] createInstanceProxyMode: req.EnvironmentId=%d, resolvedEnvironmentId=%d, mcpProtocol=%s\n", req.EnvironmentId, environmentID, mcpProtocol)

	// Create new instance record
	instance := &model.McpInstance{
		InstanceID:      instanceID,
		InstanceName:    req.Name,
		AccessType:      accessType,
		McpProtocol:     mcpProtocol,
		SourceType:      sourceType,
		SourceConfig:    json.RawMessage([]byte(req.McpServers)),
		Status:          model.InstanceStatusActive,
		IconPath:        req.IconPath,         // Add iconPath field handling
		Notes:           req.Notes,            // Add notes field handling
		McpServerID:     req.McpServerId,      // Add mcpServerId field handling
		TemplateID:      uint(req.TemplateId), // Add templateId field handling
		EnabledToken:    req.EnabledToken,
		PublicProxyPath: publicProxyPath,
		ProxyProtocol:   proxyProtocol,
		EnvironmentID:   uint(environmentID),
	}
	if req.Headers != nil && len(req.Headers) > 0 {
		if b, err := json.Marshal(req.Headers); err == nil {
			instance.Headers = b
		}
	}

	// 如果指定了环境，我们需要启动翻译代理 Sidecar，以此支持完整的代理鉴权及多重协议（如 stdio / http）的转发
	if environmentID > 0 {
		options, err := GContainerBiz.BuildProxySidecarOptions(ctx, instanceID, validationResult.Url)
		if err != nil {
			return nil, fmt.Errorf("failed to build proxy sidecar options: %w", err)
		}
		err = GContainerBiz.CreateContainer(ctx, options, environmentID, req.StartupTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to create proxy sidecar container: %w", err)
		}
		instance.ContainerName = options.ContainerName
		instance.ContainerServiceName = options.ServiceName
		instance.ContainerServiceURL = biz.getProxyContainerServiceURL(ctx, uint(environmentID), options.ContainerName, options.ServiceName)
		instance.EnvironmentID = uint(environmentID)
		// 路由协议转换为 SSE (经 sidecar 重写)
		instance.ProxyProtocol = model.McpProtocolSSE
		instance.PublicProxyPath = biz.CreatePublicProxyPath(instanceID, model.McpProtocolSSE)
	}

	tokensToSave := biz.buildCreateTokens(instanceID, req.EnabledToken, req.Tokens)
	if len(tokensToSave) > 0 {
		if err := biz.SaveTokensForInstance(ctx, tokensToSave); err != nil {
			return nil, fmt.Errorf("failed to save tokens for instance: %w", err)
		}
	}

	// Save instance to database
	if err := biz.CreateInstanceRecord(ctx, instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return &instancepb.CreateResp{
		InstanceId:  instanceID,
		Name:        req.Name,
		Status:      string(model.InstanceStatusActive),
		AccessType:  req.AccessType,
		McpProtocol: req.McpProtocol,
	}, nil
}

// createInstanceHosting Hosting mode handler function
func (biz *InstanceBiz) createInstanceHosting(ctx context.Context, req *instancepb.CreateRequest, instanceID string) (*instancepb.CreateResp, error) {
	// Validate timeout parameters
	if err := biz.validateTimeoutParams(int(req.StartupTimeout), int(req.RunningTimeout)); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %w", err)
	}
	mcpProtocol, err := common.ConvertToModelMcpProtocol(req.McpProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
	}
	sourceType, err := common.ConvertToModelSourceType(req.SourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert source type: %w", err)
	}

	if req.Port <= 0 {
		return nil, fmt.Errorf("missing required field: port")
	}
	// Validate environment ID
	if req.EnvironmentId == 0 {
		return nil, fmt.Errorf("hosting type instance requires environment ID")
	}
	if mcpProtocol == model.McpProtocolStdio {
		mcpServers := req.McpServers
		if len(mcpServers) == 0 {
			return nil, fmt.Errorf("mcp servers config is empty")
		}
		reqMcpResult, err2 := utils.ValidateMcpConfig([]byte(mcpServers))
		if err2 != nil {
			return nil, fmt.Errorf("failed to validate mcp servers: %w", err2)
		}
		if !reqMcpResult.IsValid {
			return nil, fmt.Errorf("mcp servers config is invalid: %s", reqMcpResult.ErrorMessage)
		}
		if !reqMcpResult.HasCommand {
			return nil, fmt.Errorf("mcp servers config is invalid: command is required")
		}
	}
	imageAddress := common.GetMcpHostingImage()
	containerOptions, err := GContainerBiz.BuildContainerOptions(ctx, instanceID, mcpProtocol,
		req.McpServers, req.PackageId, req.Port, req.InitScript, req.Command, imageAddress,
		req.EnvironmentVariables, req.VolumeMounts, int32(req.StartupTimeout), int32(req.RunningTimeout))
	if err != nil {
		return nil, fmt.Errorf("failed to build container options: %w", err)
	}
	err = GContainerBiz.CreateContainer(ctx, containerOptions, req.EnvironmentId, req.StartupTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Create new instance record
	containerCreateOptions, err := common.MarshalAndAssignConfig(containerOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal container create options: %w", err)
	}
	evs, err := common.MarshalAndAssignConfig(req.EnvironmentVariables)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal environment variables: %w", err)
	}
	vms, err := common.MarshalAndAssignConfig(req.VolumeMounts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal volume mounts: %w", err)
	}
	// Create target configuration
	proxyProtocol := mcpProtocol
	publicProxyPath := ""
	containerServiceURL := ""
	switch mcpProtocol {
	case model.McpProtocolStdio:
		proxyProtocol = model.McpProtocolStreamableHttp
		publicProxyPath = biz.CreatePublicProxyPath(instanceID, proxyProtocol)
		containerServiceURL = biz.getContainerServiceURL(containerOptions.ContainerName, containerOptions.ServiceName)
	case model.McpProtocolSSE, model.McpProtocolStreamableHttp:
		proxyProtocol = mcpProtocol
		publicProxyPath = biz.CreatePublicProxyPath(instanceID, proxyProtocol)
		containerServiceURL = biz.getContainerServiceURL(containerOptions.ContainerName, containerOptions.ServiceName)
	default:
		return nil, fmt.Errorf("unsupported mcp protocol: %v", mcpProtocol)
	}
	instance := &model.McpInstance{
		InstanceID:             instanceID,
		InstanceName:           req.Name,
		AccessType:             model.AccessTypeHosting,
		McpProtocol:            mcpProtocol,
		Status:                 model.InstanceStatusActive,
		PackageID:              req.PackageId,
		ContainerStatus:        model.ContainerStatusPending,
		EnvironmentID:          uint(req.EnvironmentId),
		SourceType:             sourceType,
		McpServerID:            req.McpServerId,
		TemplateID:             uint(req.TemplateId),
		EnabledToken:           req.EnabledToken,
		ImgAddr:                imageAddress,
		Port:                   req.Port,
		InitScript:             req.InitScript,
		Command:                req.Command,
		EnvironmentVariables:   evs,
		VolumeMounts:           vms,
		ContainerName:          containerOptions.ContainerName,
		ContainerServiceName:   containerOptions.ServiceName,
		ContainerIsReady:       false,
		ContainerCreateOptions: containerCreateOptions,
		ContainerLastMessage:   "container is pending",
		ContainerServiceURL:    containerServiceURL,
		StartupTimeout:         int64(req.StartupTimeout),
		RunningTimeout:         int64(req.RunningTimeout),
		SourceConfig:           json.RawMessage(req.McpServers),
		ServicePath:            req.ServicePath,
		Notes:                  req.Notes,
		IconPath:               req.IconPath,
		PublicProxyPath:        publicProxyPath,
		ProxyProtocol:          proxyProtocol,
	}
	if req.Headers != nil && len(req.Headers) > 0 {
		if b, err := json.Marshal(req.Headers); err == nil {
			instance.Headers = b
		}
	}
	tokensToSave := biz.buildCreateTokens(instanceID, req.EnabledToken, req.Tokens)
	if len(tokensToSave) > 0 {
		if err := biz.SaveTokensForInstance(ctx, tokensToSave); err != nil {
			return nil, fmt.Errorf("failed to save tokens for instance: %w", err)
		}
	}
	// Save instance to database
	if err := biz.CreateInstanceRecord(ctx, instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return &instancepb.CreateResp{
		InstanceId:  instanceID,
		Name:        req.Name,
		Status:      string(model.InstanceStatusActive),
		AccessType:  req.AccessType,
		McpProtocol: req.McpProtocol,
	}, nil
}

// validateTimeoutParams validates timeout parameters
func (biz *InstanceBiz) validateTimeoutParams(startupTimeout, runningTimeout int) error {
	// Startup timeout validation
	if startupTimeout < 0 {
		return fmt.Errorf("startup timeout cannot be negative")
	}
	if startupTimeout > 0 && startupTimeout < 30 {
		return fmt.Errorf("startup timeout cannot be less than 30 seconds")
	}
	if startupTimeout > 3600 {
		return fmt.Errorf("startup timeout cannot exceed 3600 seconds")
	}

	// Running timeout validation
	if runningTimeout < 0 {
		return fmt.Errorf("running timeout cannot be negative")
	}
	if runningTimeout > 0 && runningTimeout < 60 {
		return fmt.Errorf("running timeout cannot be less than 60 seconds")
	}
	if runningTimeout > 86400 {
		return fmt.Errorf("running timeout cannot exceed 86400 seconds")
	}

	return nil
}

// GetStatus retrieves the status of an instance
func (biz *InstanceBiz) GetStatus(ctx context.Context, req *instancepb.GetStatusRequest) (*instancepb.GetStatusResp, error) {
	instance, err := biz.GetInstance(req.InstanceId)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance information: %v", err)
	}
	if instance == nil {
		return nil, fmt.Errorf("instance does not exist")
	}

	var response *instancepb.GetStatusResp

	// Use different status query strategies based on access type
	switch instance.AccessType {
	case model.AccessTypeHosting:
		// Hosting type: query container status
		params := ContainerStatusParams{
			InstanceID: req.InstanceId,
		}
		result, err := GContainerBiz.GetContainerStatus(params)
		if err != nil {
			return nil, fmt.Errorf("failed to get container status: %s", err.Error())
		}

		response = result
	case model.AccessTypeProxy, model.AccessTypeDirect:
		_, _, mcpConfig, err := instance.GetSourceConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get source config: %v", err)
		}
		// Use HTTP probe to check service availability
		probeResult := utils.ProbePortFromURL(ctx, mcpConfig.URL, 5*time.Second)

		// Build response
		response = &instancepb.GetStatusResp{
			InstanceId: req.InstanceId,
			Status:     string(instance.Status),
		}

		// If probe fails, add error message
		if !probeResult.Success {
			response.ProbeHttp = false
		} else {
			response.ProbeHttp = true
		}
	default:
		return nil, fmt.Errorf("unsupported access type")
	}

	return response, nil
}

// GetLogs get instance logs
func (biz *InstanceBiz) GetLogs(ctx context.Context, req *instancepb.LogsRequest) (*instancepb.LogsResp, error) {
	// Set default number of lines
	lines := req.Lines
	if lines <= 0 {
		lines = 100
	}

	instance, err := biz.GetInstance(req.InstanceId)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance information: %v", err)
	}
	if instance == nil {
		return nil, fmt.Errorf("instance does not exist")
	}

	var response instancepb.LogsResp
	response.InstanceId = req.InstanceId

	// Check if it is a managed instance
	if instance.AccessType != model.AccessTypeHosting {
		response.IsManaged = false
		response.Message = "Instance is not of managed type"
		return &response, nil
	}

	response.IsManaged = true

	// Get container logs
	logs, err := GContainerBiz.GetContainerLogs(ctx, ContainerLogsParams{
		InstanceID: req.InstanceId,
		Lines:      int64(lines),
	})
	if err != nil {
		response.Message = fmt.Sprintf("Failed to get container logs: %v", err)
		return &response, nil
	}

	response.Logs = logs
	response.Message = "Logs retrieved successfully"

	return &response, nil
}

// RestartInstance restarts an instance
func (biz *InstanceBiz) RestartInstance(ctx context.Context, req *instancepb.RestartRequest) (*instancepb.RestartResp, error) {
	// 1. Query instance data by ID
	instance, err := biz.GetInstance(req.InstanceId)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, fmt.Errorf("instance does not exist")
	}

	if instance.AccessType == model.AccessTypeHosting {
		// 对 openapi 实例：重新生成 ContainerOptions 并持久化到 DB
		// 确保新代码（去掉 --header bug、更新 MCP_TARGET_URL=/mcp）生效
		if instance.SourceType == model.SourceTypeOpenapi {
			var headers map[string]string
			if len(instance.Headers) > 0 {
				_ = json.Unmarshal(instance.Headers, &headers)
			}
			newOptions, buildErr := GContainerBiz.BuildOpenapiContainerOptions(
				ctx, instance.InstanceID, instance.PackageID,
				instance.Port, int32(instance.StartupTimeout), int32(instance.RunningTimeout),
				instance.OpenapiBaseUrl, headers)
			if buildErr == nil {
				// 持久化新的 ContainerCreateOptions 到 DB（更新旧 JSON blob）
				if optJSON, marshalErr := common.MarshalAndAssignConfig(newOptions); marshalErr == nil {
					instance.ContainerCreateOptions = optJSON
					instance.ContainerServiceURL = biz.getContainerServiceURL(newOptions.ContainerName, newOptions.ServiceName)
					_ = mysql.McpInstanceRepo.Update(ctx, instance)
				}
			}
		}

		_, err = GContainerBiz.RestartContainer(ctx, instance)
		if err != nil {
			return nil, err
		}
	}
	if instance.AccessType == model.AccessTypeDirect {
		return nil, fmt.Errorf("this service does not need to be restarted")
	}

	// 3. Update container status to pending
	if err = biz.UpdateInstanceStatusToPending(ctx, instance); err != nil {
		return nil, err
	}

	pbAccessType, err := common.ConvertToProtoAccessType(instance.AccessType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access type: %w", err)
	}

	// 4. Return restart result
	return &instancepb.RestartResp{
		InstanceId: instance.InstanceID,
		Name:       instance.InstanceName,
		Status:     string(instance.Status),
		AccessType: pbAccessType,
		Message:    "Instance restarted successfully",
	}, nil
}

// GetInstance get instance information
func (biz *InstanceBiz) GetInstance(instanceID string) (*model.McpInstance, error) {
	return mysql.McpInstanceRepo.FindByInstanceID(biz.ctx, instanceID)
}

// DisableInstance disable instance
func (biz *InstanceBiz) DisableInstance(ctx context.Context, instanceID string) (string, error) {
	instance, err := biz.GetInstance(instanceID)
	if err != nil {
		return "", err
	}
	msg := "Instance has been disabled"
	if instance.AccessType == model.AccessTypeHosting {
		res, err := GContainerBiz.DeleteContainer(instance)
		if err != nil {
			return "", err
		}
		msg = res.Message
	}
	if instance.AccessType == model.AccessTypeDirect {
		return "", fmt.Errorf("this service does not need to be disabled")
	}

	instance.Status = model.InstanceStatusInactive
	instance.ContainerIsReady = false
	instance.ContainerStatus = model.ContainerStatusManualStop
	instance.ContainerLastMessage = msg
	return msg, mysql.McpInstanceRepo.Update(biz.ctx, instance)
}

func (biz *InstanceBiz) UpdateInstance(instance *model.McpInstance) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}
	err := mysql.McpInstanceRepo.Update(biz.ctx, instance)
	if err != nil {
		return err
	}
	return biz.UpdateInstanceCache(instance.InstanceID, instance)
}

// DeleteInstance deletes an instance
func (biz *InstanceBiz) DeleteInstance(instanceID string) (*instancepb.DeleteResp, error) {
	// Get instance information directly
	instance, err := mysql.McpInstanceRepo.FindByInstanceID(biz.ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance information: %w", err)
	}
	if instance == nil {
		return nil, fmt.Errorf("instance does not exist")
	}

	switch instance.AccessType {
	case model.AccessTypeHosting:
		_, err = GContainerBiz.DeleteContainer(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to delete container: %w", err)
		}
	}

	// Disable the instance and set deletion time
	err = mysql.McpInstanceRepo.Delete(biz.ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete instance: %w", err)
	}

	// Remove from cache
	cache := redis.GetMcpInstanceCache()
	cache.ClearCache(instanceID)

	return &instancepb.DeleteResp{Message: "Instance deleted successfully"}, nil
}

// updateInstance updates instance cache in Redis using CacheInstanceInfo.
// It generates a cache key and stores the instance with a default expiration.
func (biz *InstanceBiz) UpdateInstanceCache(instanceID string, instance *model.McpInstance) error {
	if instance == nil {
		return fmt.Errorf("cache instance info is nil")
	}
	cache := redis.GetMcpInstanceCache()
	key := cache.GenerateCacheKey(instanceID)
	return cache.SetRedisCacheInstance(key, instance, redis.InstanceCacheExpire)
}

// TokenListByInstanceID lists tokens for an instance with optional filters
func (biz *InstanceBiz) TokenListByInstanceID(req *instancepb.TokenListByInstanceIDRequest) (*instancepb.TokenListByInstanceIDResponse, error) {
	if req.InstanceId == "" {
		return nil, fmt.Errorf("missing required field: instanceId")
	}

	rows, err := mysql.McpTokenRepo.ListByInstanceID(biz.ctx, req.InstanceId)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %v", err)
	}

	list := make([]*instancepb.McpToken, 0, len(rows))
	for _, r := range rows {
		if req.Token != "" && req.Token != r.Token {
			continue
		}
		var usages []string
		_ = json.Unmarshal(r.Usages, &usages)
		if len(req.Usages) > 0 {
			if !usageIntersect(usages, req.Usages) {
				continue
			}
		}
		list = append(list, &instancepb.McpToken{
			Id:         int64(r.ID),
			InstanceId: r.InstanceID,
			Token:      r.Token,
			ExpireAt:   r.ExpireAt,
			PublishAt:  r.PublishAt,
			Usages:     usages,
			Enabled:    r.Enabled,
		})
	}

	return &instancepb.TokenListByInstanceIDResponse{List: list}, nil
}

func (biz *InstanceBiz) DeleteTokenByID(ctx context.Context, id int64) error {
	row, err := mysql.McpTokenRepo.FindByID(ctx, uint(id))
	if err != nil {
		return fmt.Errorf("failed to find token by id: %w", err)
	}
	if row == nil {
		return fmt.Errorf("token not found")
	}
	if err := mysql.McpTokenRepo.DeleteByID(ctx, uint(id)); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	redis.GetMcpTokenCache().Clear(row.InstanceID, row.Token)
	return nil
}

func usageIntersect(a []string, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}
	set := make(map[string]struct{}, len(a))
	for _, x := range a {
		set[x] = struct{}{}
	}
	for _, y := range b {
		if _, ok := set[y]; ok {
			return true
		}
	}
	return false
}

// SaveTokensForInstance creates or updates tokens based on incoming id field
func (biz *InstanceBiz) SaveTokensForInstance(ctx context.Context, tokens []*instancepb.McpToken) error {
	if len(tokens) == 0 {
		return nil
	}
	rows := make([]model.McpToken, 0, len(tokens))
	nowMs := time.Now().UnixMilli()
	for _, t := range tokens {
		if t.InstanceId == "" {
			return fmt.Errorf("missing required field: instanceId")
		}
		usagesBytes, _ := json.Marshal(t.Usages)
		publishAt := t.PublishAt
		if publishAt == 0 {
			publishAt = nowMs
		}

		var existing *model.McpToken
		var err error

		// 1. Try to find by ID if provided
		if t.Id > 0 {
			existing, err = mysql.McpTokenRepo.FindByID(ctx, uint(t.Id))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to find token by id: %v", err)
			}
		}

		// 2. If not found by ID (or ID is 0), try to find by Token value
		if existing == nil {
			existing, err = mysql.McpTokenRepo.FindByToken(ctx, t.InstanceId, t.Token)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to find token by value: %v", err)
			}
		}

		// 3. Update if found, otherwise prepare for creation
		if existing != nil {
			existing.InstanceID = t.InstanceId
			existing.Token = t.Token
			existing.Enabled = t.Enabled
			existing.Usages = json.RawMessage(usagesBytes)
			existing.ExpireAt = t.ExpireAt
			existing.PublishAt = publishAt
			if err := mysql.McpTokenRepo.Update(ctx, existing); err != nil {
				return fmt.Errorf("failed to update token: %v", err)
			}
			redis.GetMcpTokenCache().Clear(existing.InstanceID, existing.Token)
		} else {
			rows = append(rows, model.McpToken{
				InstanceID: t.InstanceId,
				Token:      t.Token,
				Enabled:    t.Enabled,
				Usages:     json.RawMessage(usagesBytes),
				ExpireAt:   t.ExpireAt,
				PublishAt:  publishAt,
			})
			redis.GetMcpTokenCache().Clear(t.InstanceId, t.Token)
		}
	}
	if len(rows) > 0 {
		return mysql.McpTokenRepo.CreateBatch(ctx, rows)
	}
	return nil
}

// buildCreateTokens ensures instance tokens are ready for persistence.
// When token auth is enabled and caller didn't provide tokens, create one default gateway token.
func (biz *InstanceBiz) buildCreateTokens(instanceID string, enabledToken bool, reqTokens []*instancepb.McpToken) []*instancepb.McpToken {
	nowMs := time.Now().UnixMilli()
	tokens := make([]*instancepb.McpToken, 0, len(reqTokens)+1)
	for _, token := range reqTokens {
		if token == nil {
			continue
		}
		tokenCopy := *token
		tokenCopy.InstanceId = instanceID
		if tokenCopy.PublishAt == 0 {
			tokenCopy.PublishAt = nowMs
		}
		tokens = append(tokens, &tokenCopy)
	}
	if enabledToken && len(tokens) == 0 {
		tokens = append(tokens, &instancepb.McpToken{
			InstanceId: instanceID,
			Token:      "Bearer " + uuid.NewString(),
			Enabled:    true,
			PublishAt:  nowMs,
			Usages:     []string{"default"},
		})
	}
	return tokens
}

// ListInstance get instance list
func (biz *InstanceBiz) ListInstance(page, pageSize int32, filters map[string]interface{}, sortBy, sortOrder string) (*instancepb.ListResp, error) {
	// Query data
	instances, total, err := mysql.McpInstanceRepo.FindWithPagination(biz.ctx, page, pageSize, filters, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to query instance list: %v", err)
	}

	// envIds
	envIds := make([]string, 0, len(instances))
	for _, instance := range instances {
		envIds = append(envIds, fmt.Sprintf("%d", instance.EnvironmentID))
	}
	envNames, err := mysql.McpEnvironmentRepo.FindNamesByIDs(biz.ctx, envIds)
	if err != nil {
		return nil, fmt.Errorf("failed to query environment names: %v", err)
	}

	// Convert to proto response
	instanceInfos := make([]*instancepb.ListResp_InstanceInfo, 0, len(instances))
	for _, instance := range instances {
		instanceInfo := common.ConvertToInstanceInfo(instance)
		if envName, ok := envNames[fmt.Sprintf("%d", instance.EnvironmentID)]; ok {
			instanceInfo.EnvironmentName = envName
		}
		instanceInfos = append(instanceInfos, instanceInfo)
	}

	return &instancepb.ListResp{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     instanceInfos,
	}, nil
}

// UpdateInstanceForDirect update instance
func (biz *InstanceBiz) UpdateInstanceForDirect(ctx context.Context, req *instancepb.EditRequest, oriInstance *model.McpInstance) (*instancepb.EditResp, error) {
	// Update basic information
	if req.Name != "" {
		oriInstance.InstanceName = req.Name
	}
	if req.Notes != "" {
		oriInstance.Notes = req.Notes
	}
	oriInstance.IconPath = req.IconPath

	if req.Headers != nil && len(req.Headers) > 0 {
		if b, err := json.Marshal(req.Headers); err == nil {
			oriInstance.Headers = b
		}
	} else if len(req.Headers) == 0 && req.Headers != nil {
		oriInstance.Headers = nil // allowed to clear headers
	}

	// Validate MCP configuration format
	reqMcpResult, err := utils.ValidateMcpConfig([]byte(req.McpServers))
	if err != nil {
		return nil, fmt.Errorf("failed to validate mcp servers: %w", err)
	}
	if !reqMcpResult.IsValid {
		return nil, fmt.Errorf("mcp servers config is invalid: %s", reqMcpResult.ErrorMessage)
	}
	if reqMcpResult.Url == "" {
		return nil, fmt.Errorf("mcp servers config is invalid: url is empty")
	}

	oriMcpResult, err := utils.ValidateMcpConfig([]byte(oriInstance.SourceConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to validate mcp servers: %w", err)
	}
	if !utils.CompareMcpValidationResult(reqMcpResult, oriMcpResult) {
		sourceConfig := json.RawMessage([]byte(req.McpServers))
		oriInstance.SourceConfig = sourceConfig
	}

	// Save to database
	err = mysql.McpInstanceRepo.Update(ctx, oriInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to update instance: %v", err)
	}
	// Update instance cache
	biz.UpdateInstanceCache(oriInstance.InstanceID, oriInstance)

	accessType, err := common.ConvertToProtoAccessType(oriInstance.AccessType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access type: %w", err)
	}
	mcpProtocol, err := common.ConvertToProtoMcpProtocol(oriInstance.McpProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
	}

	resp := &instancepb.EditResp{
		InstanceId:  oriInstance.InstanceID,
		Name:        oriInstance.InstanceName,
		AccessType:  accessType,
		McpProtocol: mcpProtocol,
		Status:      string(model.InstanceStatusActive),
	}

	return resp, nil
}

// UpdateInstanceForProxy update instance
func (biz *InstanceBiz) UpdateInstanceForProxy(ctx context.Context, req *instancepb.EditRequest, oriInstance *model.McpInstance) (*instancepb.EditResp, error) {
	// Update basic information
	if req.Name != "" {
		oriInstance.InstanceName = req.Name
	}
	if req.Notes != "" {
		oriInstance.Notes = req.Notes
	}
	oriInstance.IconPath = req.IconPath

	if req.Headers != nil && len(req.Headers) > 0 {
		if b, err := json.Marshal(req.Headers); err == nil {
			oriInstance.Headers = b
		}
	} else if len(req.Headers) == 0 && req.Headers != nil {
		oriInstance.Headers = nil
	}

	// Validate MCP configuration format
	reqMcpResult, err := utils.ValidateMcpConfig([]byte(req.McpServers))
	if err != nil {
		return nil, fmt.Errorf("failed to validate mcp servers: %w", err)
	}
	if !reqMcpResult.IsValid {
		return nil, fmt.Errorf("mcp servers config is invalid: %s", reqMcpResult.ErrorMessage)
	}
	if reqMcpResult.Url == "" {
		return nil, fmt.Errorf("mcp servers config is invalid: url is empty")
	}

	oriMcpResult, err := utils.ValidateMcpConfig([]byte(oriInstance.SourceConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to validate mcp servers: %w", err)
	}
	if !utils.CompareMcpValidationResult(reqMcpResult, oriMcpResult) {
		sourceConfig := json.RawMessage([]byte(req.McpServers))
		oriInstance.SourceConfig = sourceConfig
		oriInstance.ProxyProtocol = model.McpProtocol(reqMcpResult.ProtocolType)
		oriInstance.PublicProxyPath = biz.CreatePublicProxyPath(oriInstance.InstanceID, oriInstance.ProxyProtocol)

		// 如果实例配置了运行环境，当发生上游基础信息变更时，需要同步重建对应的代理 Sidecar 容器
		if oriInstance.EnvironmentID > 0 {
			options, errBuild := GContainerBiz.BuildProxySidecarOptions(ctx, oriInstance.InstanceID, reqMcpResult.Url)
			if errBuild != nil {
				return nil, fmt.Errorf("failed to build proxy sidecar options: %w", errBuild)
			}
			if optJSON, errMarshal := common.MarshalAndAssignConfig(options); errMarshal == nil {
				oriInstance.ContainerCreateOptions = optJSON
				oriInstance.ContainerName = options.ContainerName
				oriInstance.ContainerServiceName = options.ServiceName
				oriInstance.ContainerServiceURL = biz.getProxyContainerServiceURL(ctx, oriInstance.EnvironmentID, options.ContainerName, options.ServiceName)

				// 修复：K8s 环境下使用 RestartContainer 进行异步删除与重建等待，防止 "object is being deleted"
				if _, errRestart := GContainerBiz.RestartContainer(ctx, oriInstance); errRestart != nil {
					return nil, fmt.Errorf("failed to recreate proxy sidecar container: %w", errRestart)
				}
			} else {
				return nil, fmt.Errorf("failed to marshal proxy sidecar options: %w", errMarshal)
			}
		}
	}

	// Save to database
	err = mysql.McpInstanceRepo.Update(ctx, oriInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to update instance: %v", err)
	}
	// Update instance cache
	biz.UpdateInstanceCache(oriInstance.InstanceID, oriInstance)

	accessType, err := common.ConvertToProtoAccessType(oriInstance.AccessType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access type: %w", err)
	}
	mcpProtocol, err := common.ConvertToProtoMcpProtocol(oriInstance.McpProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
	}
	resp := &instancepb.EditResp{
		InstanceId:  oriInstance.InstanceID,
		Name:        oriInstance.InstanceName,
		AccessType:  accessType,
		McpProtocol: mcpProtocol,
		Status:      string(model.InstanceStatusActive),
	}

	return resp, nil
}

func (biz *InstanceBiz) UpdateInstanceForOpenapi(ctx context.Context, req *instancepb.UpdateOpenapiRequest, oriInstance *model.McpInstance) (*instancepb.EditResp, error) {
	_, err := mysql.McpOpenapiPackageRepo.FindByOpenapiFileID(ctx, req.OpenapiFileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get openapi file information: %s", err)
	}
	chooseOpenapiFileInfo, err := mysql.McpOpenapiPackageRepo.FindByOpenapiFileID(ctx, req.ChooseOpenapiFileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get choose openapi file information: %s", err)
	}
	if chooseOpenapiFileInfo.BaseOpenapiFileID != req.OpenapiFileID {
		return nil, fmt.Errorf("failed to get openapi file information")
	}

	// 判断 API 相关字段是否变化（URL、headers、openapi 文件）
	// 只有这些字段变化时才需要删除并重建容器
	oldHeadersBytes, _ := json.Marshal(oriInstance.Headers)
	newHeadersBytes, _ := json.Marshal(req.Headers)
	apiFieldChanged := req.OpenapiBaseUrl != oriInstance.OpenapiBaseUrl ||
		string(oldHeadersBytes) != string(newHeadersBytes) ||
		req.ChooseOpenapiFileID != oriInstance.PackageID

	// 更新通用字段（name/notes/icon/url/headers）到内存
	oriInstance.InstanceName = req.Name
	oriInstance.Notes = req.Notes
	oriInstance.IconPath = req.IconPath
	oriInstance.OpenapiBaseUrl = req.OpenapiBaseUrl
	if len(req.Headers) > 0 {
		if b, err2 := json.Marshal(req.Headers); err2 == nil {
			oriInstance.Headers = b
		}
	} else if req.Headers != nil {
		oriInstance.Headers = nil
	}
	if req.ChooseOpenapiFileID != oriInstance.PackageID {
		oriInstance.PackageID = req.ChooseOpenapiFileID
	}

	if !apiFieldChanged {
		// 非 API 字段变化：只保存数据库，不触发容器重建
		if err = mysql.McpInstanceRepo.Update(ctx, oriInstance); err != nil {
			return nil, fmt.Errorf("failed to update instance: %v", err)
		}
		biz.UpdateInstanceCache(oriInstance.InstanceID, oriInstance)
		return &instancepb.EditResp{
			InstanceId:  oriInstance.InstanceID,
			Name:        oriInstance.InstanceName,
			AccessType:  instancepb.AccessType_HOSTING,
			McpProtocol: instancepb.McpProtocol_STREAMABLE_HTTP,
			Status:      string(model.InstanceStatusActive),
		}, nil
	}

	// API 字段变化：重建容器
	containerOptions, err := GContainerBiz.BuildOpenapiContainerOptions(ctx, req.InstanceId, req.ChooseOpenapiFileID, oriInstance.Port, 0, 0, req.OpenapiBaseUrl, req.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to build container configuration: %v", err)
	}

	// Adjust service name for Docker runtime
	entry, err := GContainerBiz.GetRuntimeEntry(ctx, oriInstance.EnvironmentID)
	if err == nil && entry.GetRuntimeType() == container.RuntimeDocker {
		containerOptions.ServiceName = containerOptions.ContainerName
	}

	containerCreateOptions, err := common.MarshalAndAssignConfig(containerOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal container create options: %w", err)
	}

	// 更新容器状态字段
	oriInstance.ContainerCreateOptions = containerCreateOptions
	oriInstance.ContainerStatus = model.ContainerStatusPending
	oriInstance.ContainerIsReady = false
	oriInstance.ContainerLastMessage = "container is pending"

	if err = mysql.McpInstanceRepo.Update(ctx, oriInstance); err != nil {
		return nil, fmt.Errorf("failed to update instance: %v", err)
	}

	// 修复：针对 K8s 异步删除环境，统一使用 RestartContainer 解决 "object is being deleted" 竞态。
	if _, err = GContainerBiz.RestartContainer(ctx, oriInstance); err != nil {
		return nil, fmt.Errorf("failed to recreate container: %v", err)
	}

	biz.UpdateInstanceCache(oriInstance.InstanceID, oriInstance)

	return &instancepb.EditResp{
		InstanceId:  oriInstance.InstanceID,
		Name:        oriInstance.InstanceName,
		AccessType:  instancepb.AccessType_HOSTING,
		McpProtocol: instancepb.McpProtocol_STREAMABLE_HTTP,
		Status:      string(model.InstanceStatusActive),
	}, nil
}

// UpdateInstanceForHosting updates instance
func (biz *InstanceBiz) UpdateInstanceForHosting(ctx context.Context, req *instancepb.EditRequest, oriInstance *model.McpInstance) (*instancepb.EditResp, error) {
	var err error
	port := req.Port
	instanceID := req.InstanceId
	packageID := req.PackageId
	initScript := req.InitScript
	command := req.Command
	envs := req.EnvironmentVariables
	vms := req.VolumeMounts
	startupTimeout := req.StartupTimeout
	runningTimeout := req.RunningTimeout
	mcpServers := req.McpServers
	imgAddress := common.GetMcpHostingImage()

	if oriInstance.McpProtocol == model.McpProtocolStdio {
		if len(mcpServers) == 0 {
			return nil, fmt.Errorf("mcp servers config is empty")
		}
		reqMcpResult, err2 := utils.ValidateMcpConfig([]byte(mcpServers))
		if err2 != nil {
			return nil, fmt.Errorf("failed to validate mcp servers: %w", err2)
		}
		if !reqMcpResult.IsValid {
			return nil, fmt.Errorf("mcp servers config is invalid: %s", reqMcpResult.ErrorMessage)
		}
		if !reqMcpResult.HasCommand {
			return nil, fmt.Errorf("mcp servers config is invalid: command is required")
		}
		oriInstance.SourceConfig = json.RawMessage([]byte(mcpServers))
	}

	newContainerCreateOptions, err := GContainerBiz.BuildContainerOptions(ctx, instanceID, oriInstance.McpProtocol, mcpServers, packageID, port, initScript,
		command, imgAddress, envs, vms, startupTimeout, runningTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to build container configuration: %v", err)
	}

	// Check runtime type and adjust service name for Docker
	entry, err := GContainerBiz.GetRuntimeEntry(ctx, oriInstance.EnvironmentID)
	if err == nil && entry.GetRuntimeType() == container.RuntimeDocker {
		newContainerCreateOptions.ServiceName = newContainerCreateOptions.ContainerName
	}

	containerCreateOptions, err := common.MarshalAndAssignConfig(newContainerCreateOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal container create containerCreateOptions: %w", err)
	}

	// Delete old container and svc service
	_, err = GContainerBiz.DeleteContainer(oriInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to delete container: %v", err)
	}

	toEnvs, err := common.MarshalAndAssignConfig(envs)
	if err != nil {
		return nil, fmt.Errorf("failed to convert environment variables: %w", err)
	}
	toVms, err := common.MarshalAndAssignConfig(vms)
	if err != nil {
		return nil, fmt.Errorf("failed to convert volume mounts: %w", err)
	}

	// Create target configuration
	var ProxyProtocol model.McpProtocol
	publicProxyPath := ""
	containerURL := ""
	switch oriInstance.McpProtocol {
	case model.McpProtocolStdio:
		ProxyProtocol = model.McpProtocolStreamableHttp
		publicProxyPath = biz.CreatePublicProxyPath(instanceID, oriInstance.McpProtocol)
		containerURL = biz.getContainerServiceURL(newContainerCreateOptions.ContainerName, newContainerCreateOptions.ServiceName)
	case model.McpProtocolSSE, model.McpProtocolStreamableHttp:
		ProxyProtocol = oriInstance.McpProtocol
		publicProxyPath = biz.CreatePublicProxyPath(instanceID, oriInstance.McpProtocol)
		containerURL = biz.getContainerServiceURL(newContainerCreateOptions.ContainerName, newContainerCreateOptions.ServiceName)
	default:
		return nil, fmt.Errorf("unsupported mcp protocol: %v", oriInstance.McpProtocol)
	}

	// Update
	oriInstance.InstanceName = req.Name
	oriInstance.Notes = req.Notes
	oriInstance.Port = int32(port)
	oriInstance.InitScript = initScript
	oriInstance.Command = command
	oriInstance.ServicePath = req.ServicePath
	oriInstance.PackageID = packageID
	oriInstance.SourceConfig = json.RawMessage([]byte(mcpServers))
	oriInstance.ImgAddr = imgAddress
	oriInstance.EnvironmentVariables = toEnvs
	oriInstance.VolumeMounts = toVms
	oriInstance.StartupTimeout = int64(startupTimeout)
	oriInstance.RunningTimeout = int64(runningTimeout)
	oriInstance.ContainerCreateOptions = containerCreateOptions
	oriInstance.ContainerStatus = model.ContainerStatusPending
	oriInstance.ContainerServiceURL = containerURL
	oriInstance.ContainerIsReady = false
	oriInstance.PublicProxyPath = publicProxyPath
	oriInstance.ProxyProtocol = ProxyProtocol
	oriInstance.IconPath = req.IconPath

	if req.Headers != nil && len(req.Headers) > 0 {
		if b, err := json.Marshal(req.Headers); err == nil {
			oriInstance.Headers = b
		}
	} else if len(req.Headers) == 0 && req.Headers != nil {
		oriInstance.Headers = nil
	}

	// Save to database
	err = mysql.McpInstanceRepo.Update(ctx, oriInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to update instance: %v", err)
	}
	// Update instance cache
	biz.UpdateInstanceCache(oriInstance.InstanceID, oriInstance)

	accessType, err := common.ConvertToProtoAccessType(oriInstance.AccessType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access type: %w", err)
	}
	mcpProtocol, err := common.ConvertToProtoMcpProtocol(oriInstance.McpProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to convert mcp protocol: %w", err)
	}

	resp := &instancepb.EditResp{
		InstanceId:  oriInstance.InstanceID,
		Name:        oriInstance.InstanceName,
		AccessType:  accessType,
		McpProtocol: mcpProtocol,
		Status:      string(model.InstanceStatusActive),
	}
	return resp, nil
}

// GetInstancesByEnvironmentID gets instance list by environment ID
func (biz *InstanceBiz) GetInstancesByEnvironmentID(ctx context.Context, environmentID uint) ([]*model.McpInstance, error) {
	return mysql.McpInstanceRepo.FindByEnvironmentID(ctx, environmentID)
}

// CreatePublicProxyPath creates public proxy configuration
func (biz *InstanceBiz) CreatePublicProxyPath(instanceID string, mcpProtocol model.McpProtocol) string {
	addr := ""
	switch mcpProtocol {
	case model.McpProtocolSSE:
		addr = fmt.Sprintf("/%s/%s/%s", strings.Trim(common.GetGatewayRoutePrefix(), "/"), instanceID, mcpProtocol.String())
	case model.McpProtocolStreamableHttp:
		addr = fmt.Sprintf("/%s/%s/%s", strings.Trim(common.GetGatewayRoutePrefix(), "/"), instanceID, "mcp")
	default:
		addr = fmt.Sprintf("/%s/%s", strings.Trim(common.GetGatewayRoutePrefix(), "/"), instanceID)
	}
	return addr
}

// CreateInstanceRecord creates an instance record in database
func (biz *InstanceBiz) CreateInstanceRecord(ctx context.Context, instance *model.McpInstance) error {
	if instance.InstanceName == "" {
		return fmt.Errorf("instance name cannot be empty")
	}
	// Check if name already exists
	existingInstance, err := mysql.McpInstanceRepo.FindByName(biz.ctx, instance.InstanceName)
	if err == nil && existingInstance != nil {
		return fmt.Errorf("instance name %s already exists", instance.InstanceName)
	}
	// Update instance cache
	biz.UpdateInstanceCache(instance.InstanceID, instance)
	return mysql.McpInstanceRepo.Create(biz.ctx, instance)
}

// UpdateInstanceStatusToPending updates instance status to pending
func (biz *InstanceBiz) UpdateInstanceStatusToPending(ctx context.Context, instance *model.McpInstance) error {
	instance.Status = model.InstanceStatusActive
	instance.ContainerStatus = model.ContainerStatusPending
	instance.ContainerLastMessage = "Instance is restarting"
	if err := mysql.McpInstanceRepo.Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update instance status: %v", err)
	}
	return nil
}

// ListTools list tools of instance
func (biz *InstanceBiz) ListTools(ctx context.Context, instanceID string, mcpServerUrl string, token string) ([]mcp.Tool, error) {
	// 获取 mcp 客户端
	mcpClient, err := biz.getMcpClientInfo(instanceID, mcpServerUrl, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get mcp client: %s", err.Error())
	}
	defer mcpClient.Close()

	// 调用 mcp 服务的 list tools 接口
	tools, err := mcpClient.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		return nil, fmt.Errorf("list tools failed: %s", err.Error())
	}
	return tools.Tools, nil
}

// CallTool call tools of instance
func (biz *InstanceBiz) CallTool(ctx context.Context, instanceID string, toolName string, arguments any, mcpServerUrl string, token string) (interface{}, error) {
	// 获取 mcp 客户端
	mcpClient, err := biz.getMcpClientInfo(instanceID, mcpServerUrl, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get mcp client: %s", err.Error())
	}
	defer mcpClient.Close()

	// 调用 mcp 服务的 call tool 接口
	resp, err := mcpClient.CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: arguments,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("call tool failed: %s", err.Error())
	}
	return resp, nil
}

func (biz *InstanceBiz) getMcpClientInfo(instanceID string, mcpServerUrl string, token string) (*client.Client, error) {
	mcpInstance, err := biz.GetInstance(instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %s", err.Error())
	}
	if mcpInstance == nil {
		return nil, fmt.Errorf("instance does not exist")
	}

	listToolsHeader := map[string]string{}

	// 根据要求：统一严格使用传入的 token。如果外部没有传入，但是实例开启了 token 验证，则去数据库取默认 token 兼容旧的逻辑。
	if token != "" {
		listToolsHeader["Authorization"] = token
	} else if mcpInstance.EnabledToken {
		tokens, err := mysql.McpTokenRepo.ListByInstanceID(context.Background(), mcpInstance.InstanceID)
		if err != nil {
			return nil, fmt.Errorf("failed to list tokens: %s", err.Error())
		}
		var defaultTokens = getDefaultToken(tokens)
		if defaultTokens != nil {
			listToolsHeader["Authorization"] = defaultTokens.Token
		}
	}

	// Validate mcpServerUrl
	if mcpServerUrl == "" {
		return nil, fmt.Errorf("mcpServerUrl is required")
	}

	// 将前端传来的外部地址（如 http://localhost/mcp-gateway/...）转换为容器内可访问的内部地址
	// 仅对 PROXY/HOSTING 模式生效（即网关路径），DIRECT 模式的 URL 直接连接外部不做替换
	internalBase := os.Getenv("MCP_INTERNAL_BASE_URL")
	if internalBase != "" && mcpInstance.AccessType != model.AccessTypeDirect {
		// 解析 URL，替换 scheme+host 为内部地址
		if parsed, parseErr := url.Parse(mcpServerUrl); parseErr == nil {
			internalParsed, internalErr := url.Parse(internalBase)
			if internalErr == nil {
				parsed.Scheme = internalParsed.Scheme
				parsed.Host = internalParsed.Host
				mcpServerUrl = parsed.String()
				fmt.Printf("[DEBUG] getMcpClientInfo: normalized mcpServerUrl to %s\n", mcpServerUrl)
			}
		}
	}

	// Give the mcp instance its corresponding http client
	var mcpClient *client.Client
	mcpClient, err = BuildMcpClient(mcpInstance, mcpServerUrl, listToolsHeader)
	if err != nil {
		return nil, fmt.Errorf("create mcp client failed: %s", err.Error())
	}
	return mcpClient, nil
}

func BuildMcpClient(mcpInstance *model.McpInstance, mcpServerUrl string, headers map[string]string) (*client.Client, error) {
	// 1. 创建一个跳过 verify 的自定义 HTTP Client
	insecureHttpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 关键：跳过证书校验
			},
		},
	}

	var mcpClient *client.Client
	var err error

	if mcpInstance.McpProtocol == model.McpProtocolSSE {
		mcpClient, err = client.NewSSEMCPClient(
			mcpServerUrl,
			client.WithHeaders(headers),
			client.WithHTTPClient(insecureHttpClient),
		)
	} else {
		mcpClient, err = client.NewStreamableHttpClient(
			mcpServerUrl,
			transport.WithHTTPHeaders(headers),
			transport.WithHTTPBasicClient(insecureHttpClient),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("create mcp client failed: %s", err.Error())
	}

	err = mcpClient.Start(context.Background())
	if err != nil {
		return nil, fmt.Errorf("start mcp client failed: %s", err.Error())
	}

	_, err = mcpClient.Initialize(context.Background(), mcp.InitializeRequest{})
	if err != nil {
		return nil, fmt.Errorf("init mcp failed (DEBUG_TAG): %s", err.Error())
	}
	return mcpClient, nil
}

func getDefaultToken(tokens []*model.McpToken) *model.McpToken {
	var defaultTokens *model.McpToken
	for _, token := range tokens {
		if len(token.Usages) == 0 {
			continue
		}

		var usages []string
		if len(token.Usages) > 0 {
			_ = json.Unmarshal(token.Usages, &usages)
		}
		for _, usage := range usages {
			if usage == "default" {
				defaultTokens = token
				break
			}
		}
	}
	return defaultTokens
}

func (biz *InstanceBiz) getSidecarServiceName(containerName string) string {
	return containerName + common.SidecarContainerSuffix
}

// getProxyContainerServiceURL returns proxy container service URL.
// Proxy mode uses dedicated sidecar container as main runtime:
// - docker: resolve by container name
// - kubernetes: resolve by service name
func (biz *InstanceBiz) getProxyContainerServiceURL(ctx context.Context, environmentID uint, containerName, serviceName string) string {
	host := containerName
	if environmentID > 0 {
		if env, err := mysql.McpEnvironmentRepo.FindByID(ctx, environmentID); err == nil && env != nil {
			if env.Environment == model.McpEnvironmentKubernetes && serviceName != "" {
				host = serviceName
			}
		}
	}
	return fmt.Sprintf("http://%s:%d/", host, common.GetSidecarPort())
}

// resolveProxyEnvironmentID resolves runtime environment for proxy mode.
// For compatibility, when environmentId is missing (0), it falls back to default runtime environment.
func (biz *InstanceBiz) resolveProxyEnvironmentID(ctx context.Context, requestedID int32) (int32, error) {
	if requestedID > 0 {
		return requestedID, nil
	}

	defaultName := "Default-Run-Environment"
	if cfg := config.GetConfig(); cfg != nil && strings.TrimSpace(cfg.RunEnvironment.Name) != "" {
		defaultName = strings.TrimSpace(cfg.RunEnvironment.Name)
	}

	if env, err := mysql.McpEnvironmentRepo.FindByName(ctx, defaultName); err == nil && env != nil {
		return int32(env.ID), nil
	}

	environments, err := mysql.McpEnvironmentRepo.FindAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("query environments failed: %w", err)
	}
	if len(environments) == 0 {
		return 0, fmt.Errorf("no available runtime environment found")
	}

	for _, env := range environments {
		if env.Level == model.McpEnvironmentLevelSystem {
			return int32(env.ID), nil
		}
	}

	return int32(environments[0].ID), nil
}

// getContainerServiceURL 返回容器服务地址。
// 在 k8s 模式下，serviceName 与 containerName 不同（如 mcp-instance-xx-service vs mcp-instance-xx-container），
// 直接使用 serviceName（即 k8s Service DNS）。
// 在 docker 模式下，serviceName == containerName，使用 sidecar 容器名。
func (biz *InstanceBiz) getContainerServiceURL(containerName, serviceName string) string {
	var host string
	if serviceName != "" && serviceName != containerName {
		// k8s 模式：Service 名已是正确的 k8s Service DNS
		host = serviceName
	} else {
		// docker 模式：sidecar 是独立容器，用 containerName-sidecar
		host = biz.getSidecarServiceName(containerName)
	}
	return fmt.Sprintf("http://%s:%d/", host, common.GetSidecarPort())
}
