package middleware

import (
	"ginExample/common"
	"ginExample/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort()
			return
		}

		//解析token
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort()
			return
		}

		userId := claims.UserId
		db := common.GetDB()
		var user model.User
		db.First(&user, userId)
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort()
			return
		}

		//用户存在 将user信息写入上下文
		c.Set("user", user)
		c.Next()
	}
}
