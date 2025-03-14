package system_models

import "time"

type SubSystem struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;comment:子系统的唯一ID" json:"id"`
	Name           string    `gorm:"size:255;not null;comment:子系统名称" json:"name"`
	Description    string    `gorm:"type:text;comment:子系统描述" json:"description"`
	ServerAddress  string    `gorm:"size:255;not null;comment:子系统的服务器地址" json:"server_address"`
	DatabaseConfig string    `gorm:"type:text;comment:子系统的数据库配置（JSON格式）" json:"database_config"`
	CreatedAt      time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}
