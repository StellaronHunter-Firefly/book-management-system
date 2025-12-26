package middlewares

import (
	"book-management-system/utils"
	"fmt"
	"net/http"
	// "strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		fmt.Println("Authorization Header:", authHeader)
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "需要登录"})
			ctx.Abort()
			return
		}

		// tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		// fmt.Println("Extracted Token:", tokenString)
		claims, err := utils.ValidateToken(authHeader)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "令牌无效" + err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Set("userRole", claims.Role)
		ctx.Next()
	}
}
