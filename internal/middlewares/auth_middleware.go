package middlewares

import (
	"context"
	"fmt"
	"strings"
	"testcase/internal/utils"
	"testcase/package/securities"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager *securities.JWTManager
}

func NewAuthMiddleware(jwtManager *securities.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (am *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := am.extractToken(c)
		if token == "" {
			utils.ErrorResponse(c, utils.ErrUnauthorized, "Authorization token required")
			c.Abort()
			return
		}

		claims, err := am.jwtManager.ValidateAndExtract(token, "access")
		if err != nil {
			utils.ErrorResponse(c, utils.ErrUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(utils.UserIDContextKey, claims.UserID)
		c.Set(utils.UsernameContextKey, claims.Username)
		c.Set(utils.RoleContextKey, claims.Role)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, utils.UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, utils.UsernameContextKey, claims.Username)
		ctx = context.WithValue(ctx, utils.RoleContextKey, string(claims.Role))

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func (am *AuthMiddleware) AuthRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := am.extractToken(c)
		fmt.Println("Refresh Token:", token)
		if token == "" {
			utils.ErrorResponse(c, utils.ErrUnauthorized, "Refresh token required")
			c.Abort()
			return
		}

		claims, err := am.jwtManager.ValidateAndExtract(token, "refresh")
		if err != nil {
			utils.ErrorResponse(c, utils.ErrUnauthorized, "Invalid or expired refresh token")
			c.Abort()
			return
		}

		fmt.Println("Refresh Token Claims:", claims.Username)

		c.Set(utils.UserIDContextKey, claims.UserID)
		c.Set(utils.UsernameContextKey, claims.Username)
		c.Set(utils.RoleContextKey, claims.Role)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, utils.UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, utils.UsernameContextKey, claims.Username)
		ctx = context.WithValue(ctx, utils.RoleContextKey, string(claims.Role))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (am *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, utils.ErrUnauthorized, "User role not found in context")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			utils.ErrorResponse(c, utils.ErrUnauthorized, "Invalid role type")
			c.Abort()
			return
		}

		for _, requiredRole := range roles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, utils.ErrUnauthorized, "Insufficient permissions")
		c.Abort()
	}
}

func (am *AuthMiddleware) extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return authHeader
}

func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(utils.UserIDContextKey)
	if !exists {
		return 0, false
	}

	uid, ok := userID.(uint)
	return uid, ok
}

func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(utils.UserIDContextKey)
	return exists
}
