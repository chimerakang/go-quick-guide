package main

import (
	"go-todo/config"
	"go-todo/routes"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	config.LoadEnv()

	port := config.GetEnv("PORT")

	// 初始化路由
	r := routes.Init()

	// 啟動服務器
	log.Infof("Server is starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
