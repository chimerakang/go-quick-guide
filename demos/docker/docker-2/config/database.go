package config

import (
	"fmt"
	"os"

	"go-todo/models"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var log = logrus.New()

func init() {
	// 配置 logrus
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func Database() *gorm.DB {
	dbUser := GetEnv("DB_USER")
	dbPass := GetEnv("DB_PASS")
	dbHost := GetEnv("DB_HOST")
	dbPort := GetEnv("DB_PORT")
	dbName := GetEnv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	// 自動遷移 schema
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.WithError(err).Fatal("Failed to auto migrate")
	}

	// 初始化用戶
	if err := models.InitializeUsers(db); err != nil {
		log.Fatalf("Failed to initialize users: %v", err)
	}

	log.Info("Database Connection Successful and Schema Updated")

	return db
}

func LoadEnv() {
	// 檢查是否在 Docker 環境中
	if _, err := os.Stat("/.dockerenv"); err == nil {
		log.Info("Running in Docker environment, skipping .env file load")
		return
	}

	// 不在 Docker 環境中，嘗試加載 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found, using environment variables")
	} else {
		log.Info(".env file loaded successfully")
	}
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.WithField("key", key).Fatal("Required environment variable not set")
	}
	return value
}
