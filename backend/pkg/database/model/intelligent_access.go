package model

import "time"

// IntelligentAccess 接入信息
type IntelligentAccess struct {
	ID         int64  `gorm:"primaryKey;autoIncrement;column:access_id;comment:接入信息ID"`
	AccessName string `gorm:"type:varchar(255);not null;column:access_name;comment:名称"`
	AccessType string `gorm:"type:varchar(255);not null;column:access_type;comment:类型: COZE, DifyEnterprise, N8N"`

	SubType      string `gorm:"type:varchar(255);not null;column:sub_type;comment:接入类型: Team, Personal"`
	EnterpriseID string `gorm:"type:varchar(255);not null;column:enterprise_id;comment:企业ID"`
	CozeUserID   string `gorm:"type:varchar(255);not null;column:coze_user_id;comment:Coze用户ID"`

	DbHost     string `gorm:"type:varchar(255);column:db_host;comment:dify数据库地址"`
	DbPort     int    `gorm:"type:int;column:db_port;comment:dify数据库端口"`
	DbUser     string `gorm:"type:varchar(32);column:db_user;comment:dify数据库帐号"`
	DbPassword string `gorm:"type:varchar(255);column:db_password;comment:dify数据库密码"`
	DbName     string `gorm:"type:varchar(255);column:db_name;comment:dify数据库名"`

	BaseUrl  string `gorm:"type:varchar(255);column:base_url;comment:n8n基础URL"`
	Username string `gorm:"type:varchar(255);column:username;comment:n8n用户名"`
	Password string `gorm:"type:varchar(255);column:password;comment:n8n密码"`

	CreateTime time.Time `gorm:"type:datetime;column:create_time;comment:创建时间"`
	UpdateTime time.Time `gorm:"type:datetime;column:update_time;comment:更新时间"`
}

var tableName string

func SetIntelligentAccessTableName(name string) {
	tableName = name
}

func (m *IntelligentAccess) TableName() string {
	return tableName
}
