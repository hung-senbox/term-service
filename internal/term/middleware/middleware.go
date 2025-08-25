package middleware

import (
	"context"
	"net/http"
	"strings"
	"term-service/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Secured() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		if len(authorizationHeader) == 0 {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]

		token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// --- UserID ---
			if userId, ok := claims[constants.UserID.String()].(string); ok {
				// gin context → key phải là string
				c.Set(constants.UserID.String(), userId)
				// request context → key là ContextKey
				ctx := context.WithValue(c.Request.Context(), constants.UserID, userId)
				c.Request = c.Request.WithContext(ctx)
			}

			// --- UserName ---
			if userName, ok := claims[constants.UserName.String()].(string); ok {
				c.Set(constants.UserName.String(), userName)
				ctx := context.WithValue(c.Request.Context(), constants.UserName, userName)
				c.Request = c.Request.WithContext(ctx)
			}

			// --- Roles ---
			if userRoles, ok := claims[constants.UserRoles.String()].(string); ok {
				c.Set(constants.UserRoles.String(), userRoles)
				ctx := context.WithValue(c.Request.Context(), constants.UserRoles, userRoles)
				c.Request = c.Request.WithContext(ctx)
			}
		}

		// Token
		c.Set(constants.Token.String(), tokenString)
		ctx := context.WithValue(c.Request.Context(), constants.Token, tokenString)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesAny, exists := c.Get(constants.UserRoles.String())
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Roles not found"})
			return
		}

		rolesStr, ok := rolesAny.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid roles format"})
			return
		}

		// ví dụ roles: "SuperAdmin, Teacher"
		roles := strings.Split(rolesStr, ",")
		isAdmin := false
		for _, role := range roles {
			if strings.TrimSpace(role) == "SuperAdmin" {
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
