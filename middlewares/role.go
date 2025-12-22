package middlewares

import (
	"book-management-system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleValue, exists := ctx.Get("userRole")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			ctx.Abort()
			return
		}

		role, ok := roleValue.(models.UserRole)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "角色类型错误"})
			ctx.Abort()
			return
		}

		if role != models.RoleAdmin {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
