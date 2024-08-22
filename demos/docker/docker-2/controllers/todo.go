package controllers

import (
	"net/http"

	"go-todo/config"
	"go-todo/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
	log      *logrus.Logger
)

func init() {
	database = config.Database()
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// func Show(c *gin.Context) {
// 	var todos []models.Todo
// 	result := database.Find(&todos)
// 	if result.Error != nil {
// 		log.WithError(result.Error).Error("Failed to fetch todos")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 		return
// 	}

// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"Todos": todos,
// 	})
// }

// func Add(c *gin.Context) {
// 	item := c.PostForm("item")
// 	todo := models.Todo{Item: item}

// 	result := database.Create(&todo)
// 	if result.Error != nil {
// 		log.WithError(result.Error).Error("Failed to add todo")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 		return
// 	}

// 	c.Redirect(http.StatusMovedPermanently, "/")
// }

// func Delete(c *gin.Context) {
// 	id := c.Param("id")

// 	result := database.Delete(&models.Todo{}, id)
// 	if result.Error != nil {
// 		log.WithError(result.Error).Error("Failed to delete todo")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 		return
// 	}

// 	c.Redirect(http.StatusMovedPermanently, "/")
// }

// func Complete(c *gin.Context) {
// 	id := c.Param("id")

// 	result := database.Model(&models.Todo{}).Where("id = ?", id).Update("completed", true)
// 	if result.Error != nil {
// 		log.WithError(result.Error).Error("Failed to update todo")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 		return
// 	}

// 	c.Redirect(http.StatusMovedPermanently, "/")
// }

func Show(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var todos []models.Todo
	result := database.Where("user_id = ?", userID).Find(&todos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Todos": todos,
	})
}

func Add(c *gin.Context) {
	userID, _ := c.Get("user_id")
	item := c.PostForm("item")
	todo := models.Todo{Item: item, UserID: userID.(uint)}

	result := database.Create(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add todo"})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	result := database.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Todo{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func Complete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	result := database.Model(&models.Todo{}).Where("id = ? AND user_id = ?", id, userID).Update("completed", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.Redirect(http.StatusFound, "/")
}
