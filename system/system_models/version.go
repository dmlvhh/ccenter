package system_models

import "time"

type Version struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;comment:版本记录的唯一ID" json:"id"`
	SubSystemID   uint      `gorm:"not null;comment:关联的子系统ID" json:"sub_system_id"`
	VersionNumber string    `gorm:"size:50;not null;comment:版本号（如v1.0.0）" json:"version_number"`
	ReleaseNotes  string    `gorm:"type:text;comment:版本发布说明" json:"release_notes"`
	ReleaseDate   time.Time `gorm:"not null;comment:发布日期" json:"release_date"`
	CreatedAt     time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	SubSystem     SubSystem `gorm:"foreignKey:SubSystemID;comment:关联的子系统" json:"sub_system"`
}
