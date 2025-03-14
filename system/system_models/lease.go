package system_models

import "time"

type Lease struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;comment:租赁记录的唯一ID" json:"id"`
	SubSystemID uint      `gorm:"not null;comment:关联的子系统ID" json:"sub_system_id"`
	StartTime   time.Time `gorm:"not null;comment:租赁开始时间" json:"start_time"`
	EndTime     time.Time `gorm:"not null;comment:租赁结束时间" json:"end_time"`
	Status      string    `gorm:"size:50;not null;comment:租赁状态（如active, expired）" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	SubSystem   SubSystem `gorm:"foreignKey:SubSystemID;comment:关联的子系统" json:"sub_system"`
}
