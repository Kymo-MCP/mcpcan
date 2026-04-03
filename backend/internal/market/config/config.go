package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kymo-mcp/mcpcan/pkg/common"
	"github.com/kymo-mcp/mcpcan/pkg/utils"
	"github.com/kymo-mcp/mcpcan/pkg/version"

	"github.com/spf13/viper"
)

var GlobalConfig *Config

// Config represents configuration structure
type Config struct {
	// RunMode indicates current running mode, e.g., "demo" or "prod" or "kymo"
	RunMode     string                `mapstructure:"runMode"`
	ServiceName string                `mapstructure:"-"`
	VersionInfo *version.VersionInfo  `mapstructure:"-"`
	Server      common.ServerConfig   `mapstructure:"server"`
	Database    common.DatabaseConfig `mapstructure:"database"`
	Code        common.CodeConfig     `mapstructure:"code"`
	Market      common.MarketConfig   `mapstructure:"market"`
	Log         common.LogConfig      `mapstructure:"log"`
	Secret      string                `mapstructure:"secret"`
	// Storage configuration
	Storage common.StorageConfig `mapstructure:"storage"`
	// Global Domain
	Domain string `mapstructure:"domain"`
	// RunEnvironment configuration
	RunEnvironment common.RunEnvironmentConfig `mapstructure:"runEnvironment"`
	// DemoMaxInstances specifies the maximum number of active instances allowed in demo mode
	DemoMaxInstances int            `mapstructure:"demoMaxInstances"`
	Init             InitUserConfig `mapstructure:"init"`
	// CodeMode indicates whether it is OpenCode or EnterpriseCode
	CodeMode common.CodeMode `mapstructure:"-"`
}

// InitUserConfig represents admin user initialization configuration
type InitUserConfig struct {
	AdminUsername        string `mapstructure:"admin_username"`
	AdminPassword        string `mapstructure:"admin_password"`
	AdminPasswordPolicy  string `mapstructure:"admin_password_policy"`
	MenuSyncPolicy       string `mapstructure:"menu_sync_policy"`
	AdminNickname        string `mapstructure:"admin_nickname"`
	AdminRoleName        string `mapstructure:"admin_role_name"`
	AdminRoleDescription string `mapstructure:"admin_role_description"`
	AdminRoleLevel       int    `mapstructure:"admin_role_level"`
	AdminDataScope       string `mapstructure:"admin_data_scope"`
	AdminDeptName        string `mapstructure:"admin_dept_name"`
}

const (
	InitAdminPasswordPolicyCreateOnly = "create_only"
	InitAdminPasswordPolicyForceReset = "force_reset"

	InitMenuSyncPolicyOff      = "off"
	InitMenuSyncPolicyAddOnly  = "add_only"
	InitMenuSyncPolicyFullSync = "full_sync"
)

var serviceName = "market"
var cfgFileName = "market.yaml"

// GetConfig gets global configuration
func GetConfig() *Config {
	return GlobalConfig
}

// Load loads configuration file
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// If configuration file path is not specified, try to find it automatically
	var err error
	configPath, err := common.FindConfigFile(cfgFileName)
	if err != nil {
		return nil, err
	}

	// Set configuration file path
	v.SetConfigFile(configPath)

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Read run mode from environment if not provided in config
	if strings.TrimSpace(config.RunMode) == "" {
		if envMode := strings.TrimSpace(config.RunMode); envMode != "" {
			config.RunMode = envMode
		}
	}

	// Read demo max instances from environment if provided
	if envDemoMax := strings.TrimSpace(os.Getenv("DEMO_MAX_INSTANCES")); envDemoMax != "" {
		if v, err := strconv.Atoi(envDemoMax); err == nil && v > 0 {
			config.DemoMaxInstances = v
		}
	}

	if config.Code.Upload.MaxFileSize == 0 {
		config.Code.Upload.MaxFileSize = 100
	}

	if config.Code.Upload.AllowedExtensions == nil {
		config.Code.Upload.AllowedExtensions = []string{".zip", ".tar.gz", ".tar", ".rar"}
	}

	if config.Storage.RootPath == "" {
		config.Storage.RootPath = "/app/data"
	}
	utils.MkdirP(config.Storage.RootPath)

	if config.Storage.CodePath == "" {
		config.Storage.CodePath = "/app/data/code-package"
	}
	utils.MkdirP(config.Storage.CodePath)

	if config.Storage.StaticPath == "" {
		config.Storage.StaticPath = "/app/data/static"
	}
	utils.MkdirP(config.Storage.StaticPath)

	if config.Storage.OpenapiFilePath == "" {
		config.Storage.OpenapiFilePath = "/app/data/openapi-file"
	}
	utils.MkdirP(config.Storage.OpenapiFilePath)

	// Init policy defaults (safe for upgrade)
	switch strings.ToLower(strings.TrimSpace(config.Init.AdminPasswordPolicy)) {
	case InitAdminPasswordPolicyForceReset:
		config.Init.AdminPasswordPolicy = InitAdminPasswordPolicyForceReset
	default:
		config.Init.AdminPasswordPolicy = InitAdminPasswordPolicyCreateOnly
	}
	switch strings.ToLower(strings.TrimSpace(config.Init.MenuSyncPolicy)) {
	case InitMenuSyncPolicyOff:
		config.Init.MenuSyncPolicy = InitMenuSyncPolicyOff
	case InitMenuSyncPolicyFullSync:
		config.Init.MenuSyncPolicy = InitMenuSyncPolicyFullSync
	default:
		config.Init.MenuSyncPolicy = InitMenuSyncPolicyAddOnly
	}

	// Set sensible defaults
	if config.DemoMaxInstances <= 0 {
		config.DemoMaxInstances = 5
	}

	// Validate mandatory ServiceName
	if strings.TrimSpace(config.Server.ServiceName) == "" {
		return nil, fmt.Errorf("server.serviceName is required in config, but not found")
	}

	// Append Version information
	config.ServiceName = serviceName
	config.VersionInfo = version.GetVersionInfo()

	// Prioritize CodeMode from environment variable for runtime switching (Single Image strategy)
	if envCodeMode := os.Getenv("CODE_MODE"); envCodeMode != "" {
		config.CodeMode = common.CodeMode(envCodeMode)
	} else {
		config.CodeMode = common.CodeMode(config.VersionInfo.CodeMode)
	}

	// If Market.Host is empty, fallback to top-level Domain
	if strings.TrimSpace(config.Market.Host) == "" && config.Domain != "" {
		config.Market.Host = config.Domain
	}

	GlobalConfig = &config

	return &config, nil
}

// IsDemoMode returns true if the application is running in demo mode
func IsDemoMode() bool {
	if GlobalConfig == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(GlobalConfig.RunMode), "demo")
}

// GetDemoMaxInstances returns the maximum allowed active instances in demo mode
func GetDemoMaxInstances() int {
	if GlobalConfig == nil || GlobalConfig.DemoMaxInstances <= 0 {
		return 10
	}
	return GlobalConfig.DemoMaxInstances
}
