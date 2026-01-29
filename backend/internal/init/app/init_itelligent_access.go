package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
	"github.com/kymo-mcp/mcpcan/pkg/logger"
	"gorm.io/gorm"
)

type intelligentAccessSeed struct {
	AccessName   string `json:"access_name"`
	AccessType   string `json:"access_type"`
	DbHost       string `json:"db_host"`
	DbPort       int    `json:"db_port"`
	DbUser       string `json:"db_user"`
	DbPassword   string `json:"db_password"`
	DbName       string `json:"db_name"`
	SubType      string `json:"sub_type"`
	EnterpriseID string `json:"enterprise_id"`
	BaseUrl      string `json:"base_url"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	CozeUserID   string `json:"coze_user_id"`
}

const intelligentAccessSeedJSON = `[
  {
    "access_name": "COZE_SAAS_TEAM(Complete the enterprise ID information)",
    "access_type": "COZE",
    "sub_type": "Team",
    "enterprise_id": "Please fill in your enterprise ID here"
  },
  {
    "access_name": "COZE_SAAS_PERSON(Complete the Coze user ID information)",
    "access_type": "COZE",
    "sub_type": "Person",
    "coze_user_id": "Please fill in your Coze user ID here"
  },
  {
    "access_name": "DifyEnterprise_DB_INFO(Please fill in the database information)",
    "access_type": "DifyEnterprise",
    "db_host": "postgres_db_host_address",
    "db_port": 5432,
    "db_user": "postgres_db_user",
    "db_password": "postgres_db_password",
    "db_name": "dify"
  },
  {
    "access_name": "Dify_DB_INFO(Please fill in the database information)",
    "access_type": "Dify",
    "db_host": "postgres_db_host_address",
    "db_port": 5432,
    "db_user": "postgres_db_user",
    "db_password": "postgres_db_password",
    "db_name": "dify"
  },
  {
    "access_name": "N8N_Account(Please fill in the account information)",
    "access_type": "N8N",
    "base_url": "http://172.16.40.5:5678",
    "username": "test@test.com",
    "password": "pwd123456"
  }
]`

func (a *App) initIntelligentAccess(ctx context.Context) error {
	var seeds []intelligentAccessSeed
	if err := json.Unmarshal([]byte(intelligentAccessSeedJSON), &seeds); err != nil {
		return fmt.Errorf("failed to parse intelligent access seed data: %w", err)
	}

	if len(seeds) == 0 {
		return nil
	}

	db := mysql.GetDB().Model(&model.IntelligentAccess{})
	createdCount := 0
	skippedCount := 0

	now := time.Now()
	for _, seed := range seeds {
		record := &model.IntelligentAccess{
			AccessName:   seed.AccessName,
			AccessType:   seed.AccessType,
			DbHost:       seed.DbHost,
			DbPort:       seed.DbPort,
			DbUser:       seed.DbUser,
			DbPassword:   seed.DbPassword,
			DbName:       seed.DbName,
			SubType:      seed.SubType,
			EnterpriseID: seed.EnterpriseID,
			BaseUrl:      seed.BaseUrl,
			Username:     seed.Username,
			Password:     seed.Password,
			CozeUserID:   seed.CozeUserID,
			CreateTime:   now,
			UpdateTime:   now,
		}

		var existing model.IntelligentAccess
		err := db.WithContext(ctx).Where("access_name = ?", seed.AccessName).First(&existing).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to query intelligent access record: %w", err)
		}

		if err == nil {
			skippedCount++
			continue
		}

		if err := db.WithContext(ctx).Create(record).Error; err != nil {
			return fmt.Errorf("failed to create intelligent access record: %w", err)
		}
		createdCount++
	}

	logger.Info("Intelligent access initialization completed",
		zap.Int("created", createdCount),
		zap.Int("skipped", skippedCount))
	return nil
}
