package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func InitializeUsers(db *gorm.DB) error {
	// 自動遷移 User 模型（創建表如果不存在）
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	// 檢查是否已存在用戶
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return nil // 如果已存在用戶，則不需要創建默認用戶
	}

	// 創建默認用戶
	defaultUsers := []User{
		{Username: "admin", Password: "admin123"},
		{Username: "user", Password: "user123"},
		{Username: "guest", Password: "guest123"},
	}

	for i := range defaultUsers {
		hashedPassword, err := HashPassword(defaultUsers[i].Password)
		if err != nil {
			return err
		}
		defaultUsers[i].Password = hashedPassword
	}

	result := db.Create(&defaultUsers)
	return result.Error
}
