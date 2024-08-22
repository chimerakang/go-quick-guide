package controllers

import (
	"fmt"
	"net/http"

	"go-todo/controllers/jwt"
	"go-todo/models"

	"github.com/gin-gonic/gin"
)

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func ProcessLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 這裡應該有實際的用戶驗證邏輯
	var user models.User
	if err := database.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !models.CheckPasswordHash(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// if username == "admin" && password == "password" { // 示例驗證
	tokenString, err := jwt.GenerateToken(user.Username, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	// Set cookie
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	fmt.Printf("Token created: %s\n", tokenString)

	// 如果验证成功，重定向到 Todo 列表页面
	c.Redirect(http.StatusSeeOther, "/")
	// c.JSON(http.StatusOK, gin.H{"token": tokenString})
	// } else {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	// }
}

// @Summary User logout
// @Description Logout the current user
// @Produce json
// @Success 303 {string} string "See Other"
// @Router /logout [get]
func ProcessLogout(c *gin.Context) {
	// loggedInUser = ""
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
