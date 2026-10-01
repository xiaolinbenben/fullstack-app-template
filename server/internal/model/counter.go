// Package model 存放 GORM 模型。新模型要在 database.Migrate 里注册。
package model

import "time"

// WoodfishName 是敲木鱼计数行的固定名称。
const WoodfishName = "woodfish"

// Counter 保存一个具名计数。
type Counter struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:32;uniqueIndex;not null"`
	Count     int64  `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
