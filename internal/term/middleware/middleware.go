package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"term-service/internal/gateway"
	"term-service/logger"
	"term-service/pkg/constants"
	"term-service/pkg/helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Secured(userGw gateway.UserGateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		// app language header
		appLanguage := helper.ParseAppLanguage(c.GetHeader("X-App-Language"), 1)
		c.Writer.Header().Set("X-App-Language", strconv.Itoa(int(appLanguage)))
		c.Set(constants.AppLanguage.String(), appLanguage)
		ctx := context.WithValue(c.Request.Context(), constants.AppLanguage, appLanguage)
		c.Request = c.Request.WithContext(ctx)

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
		ctx = context.WithValue(c.Request.Context(), constants.Token, tokenString)
		c.Request = c.Request.WithContext(ctx)

		// --- Call user-service to get current user ---
		currentUser, err := userGw.GetCurrentUser(
			context.WithValue(ctx, constants.CurrentUserKey, tokenString),
		)
		if err != nil {
			logger.WriteLogEx("error", "failed to get current user", map[string]any{
				"error": err.Error(),
			})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		// --- Set currentUser vào context ---
		c.Set(string(constants.CurrentUserKey), currentUser)
		ctx = context.WithValue(ctx, constants.CurrentUserKey, currentUser)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func SecuredV2(userGW gateway.UserGateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		if len(authorizationHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "authorization header required",
			})
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]

		// parse unverified để extract claims nếu cần
		token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// optional: bạn vẫn có thể lấy claim trước khi call user-service
			if userId, ok := claims[constants.UserID.String()].(string); ok {
				c.Set(constants.UserID.String(), userId)
			}
		}

		// gọi user-service để lấy current user
		currentUser, err := userGW.GetCurrentUser(
			context.WithValue(c.Request.Context(), constants.CurrentUserKey, tokenString),
		)
		if err != nil {
			logger.WriteLogEx("error", "failed to get current user", map[string]any{
				"error": err.Error(),
			})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		// set currentUser vào gin.Context
		c.Set(string(constants.CurrentUserKey), currentUser)

		// set vào request.Context để service khác có thể lấy
		ctx := context.WithValue(c.Request.Context(), constants.CurrentUserKey, currentUser)
		c.Request = c.Request.WithContext(ctx)

		// cũng set token để reuse
		c.Set(constants.Token.String(), tokenString)
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), constants.Token, tokenString),
		)

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

func RequireIsSuperAdminOrOrgAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, ok := helper.CurrentUserFromCtx(c.Request.Context())
		if !ok || currentUser == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "current user not found in context",
			})
			return
		}

		// nếu user là SuperAdmin thì pass ngay
		if currentUser.IsSuperAdmin {
			c.Next()
			return
		}

		// kiểm tra Roles slice
		if currentUser.Roles == nil || len(*currentUser.Roles) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "roles not found for current user",
			})
			return
		}

		isAllowed := false
		for _, role := range *currentUser.Roles {
			if strings.EqualFold(role.RoleName, "SuperAdmin") {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "permission denied: require SuperAdmin or OrgAdmin",
			})
			return
		}

		c.Next()
	}
}
