package app

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/kymo-mcp/mcpcan/internal/init/config"
	"github.com/kymo-mcp/mcpcan/pkg/common"
	dbpkg "github.com/kymo-mcp/mcpcan/pkg/database"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
	"github.com/kymo-mcp/mcpcan/pkg/logger"
)

// App application structure
type App struct {
	config          *config.InitConfig
	logger          *zap.Logger
	adminUser       *model.SysUser
	codePackageList []*model.McpCodePackage
	mcpTemplateList []*model.McpTemplate
}

// New creates application instance
func New() *App {
	// Load configuration
	if err := config.Load(); err != nil {
		return nil
	}
	// Initialize logger
	if err := logger.Init(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Format); err != nil {
		return nil
	}

	// Print configuration information
	logger.Debug("Version info", zap.String("version", fmt.Sprintf("%+v", config.GlobalConfig.VersionInfo)))

	return &App{
		config:          config.GlobalConfig,
		logger:          logger.L().Logger,
		adminUser:       &model.SysUser{},
		codePackageList: make([]*model.McpCodePackage, 0),
		mcpTemplateList: make([]*model.McpTemplate, 0),
	}
}

// Initialize initializes the application
func (a *App) Initialize() error {
	// Initialize database
	if err := a.loadMysql(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	logger.Info("Init service initialized successfully")
	return nil
}

// loadMysql initializes MySQL database connection and loads necessary tables
func (a *App) loadMysql() error {
	// Common tables for all modes
	tableInitializers := []func() (string, error){
		func() (string, error) {
			repo := mysql.NewMcpCodePackageRepository()
			return (&model.McpCodePackage{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewMcpEnvironmentRepository()
			return (&model.McpEnvironment{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewGatewayLogRepository()
			return (&model.GatewayLog{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewMcpInstanceRepository()
			return (&model.McpInstance{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewMcpMigrationRepository()
			return (&model.Migration{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewMcpOpenapiPackageRepository()
			return (&model.McpOpenapiPackage{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewMcpTemplateRepository()
			return (&model.McpTemplate{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewMcpToIntelligentTaskRepository()
			return (&model.McpToIntelligentTask{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewMcpToIntelligentTaskLogRepository()
			return (&model.McpToIntelligentTaskLog{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			repo := mysql.NewMcpTokenRepository()
			return (&model.McpToken{}).TableName(), repo.InitTable()
		},
		func() (string, error) {
			if a.config.RunMode == common.RunModeKymo {
				model.SetIntelligentAccessTableName("intelligent_access")
				mod := &model.IntelligentAccess{}
				return mod.TableName(), nil
			}
			model.SetIntelligentAccessTableName("mcpcan_intelligent_access")
			mod := &model.IntelligentAccess{}
			repo := mysql.NewIntelligentAccessRepository()
			return mod.TableName(), repo.InitTable()
		},
	}

	// Sys tables - only load when NOT in kymo mode
	if a.config.RunMode != common.RunModeKymo {
		tableInitializers = append(tableInitializers,
			func() (string, error) {
				repo := mysql.NewSysDeptRepository()
				return (&model.SysDept{}).TableName(), repo.InitTable()
			},
			func() (string, error) {
				repo := mysql.NewSysEncryptionKeyRepository(mysql.GetDB())
				return (&model.SysEncryptionKey{}).TableName(), repo.InitTable()
			},
			func() (string, error) {
				repo := mysql.NewSysRoleRepository()
				return (&model.SysRole{}).TableName(), repo.InitTable()
			},
			func() (string, error) {
				repo := mysql.NewSysRolesDeptsRepository()
				return (&model.SysRolesDepts{}).TableName(), repo.InitTable()
			},
			func() (string, error) {
				repo := mysql.NewSysUserRepository()
				return (&model.SysUser{}).TableName(), repo.InitTable()
			},
			func() (string, error) {
				repo := mysql.NewSysUsersRolesRepository()
				return (&model.SysUsersRoles{}).TableName(), repo.InitTable()
			},
		)
	}

	return dbpkg.Init(&a.config.Database.MySQL, tableInitializers...)
}

// Run 运行应用程序
func (a *App) Run() error {
	// 拷贝项目路径下 data 目录所有基础数据到挂载根目录中
	if err := a.copyInitData("./init-data/static", a.config.Storage.StaticPath); err != nil {
		return fmt.Errorf("failed to copy data directory: %w", err)
	}

	if a.config.RunMode != common.RunModeKymo {
		// 创建管理员用户
		_, err := a.createAdminUser()
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.initDataScope(context.Background()); err != nil {
			logger.Error("Failed to init data scope", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunMigrations()
	}()

	wg.Wait()

	logger.Info("Shutting down init service...")

	// 优雅关闭
	return a.Shutdown()
}

// Shutdown 优雅关闭应用程序
func (a *App) Shutdown() error {
	// 关闭数据库连接
	if err := mysql.Close(); err != nil {
		logger.Error("Failed to close database", zap.Error(err))
		return err
	}

	logger.Info("Init service shutdown completed")
	return nil
}

// initDataScope creates the default environment
func (a *App) initDataScope(ctx context.Context) error {
	// 初始化代码包数据
	if err := a.initCodePackage(ctx); err != nil {
		return fmt.Errorf("failed to init code package data: %w", err)
	}
	// 初始化 OpenAPI 文档数据
	if err := a.initOpenapi(ctx); err != nil {
		return fmt.Errorf("failed to init openapi data: %w", err)
	}
	// 初始化智能访问数据
	if err := a.initIntelligentAccess(ctx); err != nil {
		return fmt.Errorf("failed to init intelligent access data: %w", err)
	}
	// 初始化 MCP 模板数据（使用嵌入式模板 JSON）
	if err := a.initMcpTemplateData(ctx); err != nil {
		return fmt.Errorf("failed to init mcp template data: %w", err)
	}
	return nil
}
