package middleware

import (
	// "net/http"
	// "strings"
	// "term-info-service/pkg/constants"

	"net/http"
	"strings"
	"term-service/pkg/constants"

	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func Secured() gin.HandlerFunc {
	return func(context *gin.Context) {
		authorizationHeader := context.GetHeader("Authorization")

		if len(authorizationHeader) == 0 {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]

		token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userId, ok := claims[constants.UserID].(string); ok {
				context.Set(constants.UserID, userId)
			}

			if userName, ok := claims[constants.UserName].(string); ok {
				context.Set(constants.UserName, userName)
			}

			if userRoles, ok := claims[constants.UserRoles].(string); ok {
				context.Set(constants.UserRoles, userRoles)
			}
		}

		context.Set(constants.Token, tokenString)
		context.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesAny, exists := c.Get(constants.UserRoles)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Roles not found"})
			return
		}

		rolesStr, ok := rolesAny.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid roles format"})
			return
		}

		// Chuyển chuỗi "Admin, Teacher" thành slice
		roles := strings.Split(rolesStr, ",")
		isAdmin := false
		for _, role := range roles {
			if strings.TrimSpace(role) == "Admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		c.Next()
	}
}
