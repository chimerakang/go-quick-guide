package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key") // 在實際應用中，應該從環境變量或配置文件中讀取

type CustomClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成一個新的 JWT token
func GenerateToken(username string, userId uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &CustomClaims{
		Username: username,
		UserID:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateToken 驗證 JWT token
// func ValidateToken(tokenString string) (*CustomClaims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, errors.New("invalid token")
// }

// Function to verify JWT tokens
func ValidateToken(tokenString string) (*CustomClaims, error) {

	claims := &CustomClaims{}
	// Parse the token with the secret key
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return claims, nil
}

// AuthMiddleware Gin 中間件，用於檢查 JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			fmt.Println("Token missing in cookie")
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// Verify the token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			fmt.Printf("Token verification failed: %v\\n", err)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// Extract user ID from the token claims
		userID := claims.UserID

		// Print information about the verified token
		fmt.Printf("Token verified successfully. User ID: %d\n", userID)

		// Set the user ID in the context for use in subsequent handlers
		c.Set("user_id", userID)

		// Continue with the next middleware or route handler
		c.Next()
	}

}
