package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"  structs:"-"`                     // 主键ID
	CreatedAt time.Time      `gorm:"index;type:datetime(0)" json:"created_at"  structs:"-"` // 创建时间
	UpdatedAt time.Time      `gorm:"index;type:datetime(0)" json:"updated_at"  structs:"-"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"type:datetime(0);index" json:"-" structs:"-"`           // 删除时间
}

// Pagination 是分页查询的基本参数。
type Pagination struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Sort     string `form:"sort" json:"sort"`
	Key      string `form:"key" json:"key"`
}
