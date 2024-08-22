package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	ID        uint   `gorm:"primaryKey"`
	Item      string `gorm:"type:text;not null"`
	Completed bool   `gorm:"default:false"`
	UserID    uint   // 外鍵
	User      User   `gorm:"foreignKey:UserID"` // 關聯
}

// 為了向後兼容，我們可以添加一個方法來獲取 ID
func (t *Todo) GetID() uint {
	return t.ID
}

// 初始化 Todo 表
func InitializeTodos(db *gorm.DB) error {
	return db.AutoMigrate(&Todo{})
}
