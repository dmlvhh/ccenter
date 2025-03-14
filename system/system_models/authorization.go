package system_models

import "time"

type Authorization struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;comment:授权记录的唯一ID" json:"id"`
	SubSystemID uint      `gorm:"not null;comment:关联的子系统ID" json:"sub_system_id"`
	Token       string    `gorm:"size:255;not null;comment:授权令牌（如JWT）" json:"token"`
	ExpiresAt   time.Time `gorm:"not null;comment:授权过期时间" json:"expires_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	SubSystem   SubSystem `gorm:"foreignKey:SubSystemID;comment:关联的子系统" json:"sub_system"`
}
