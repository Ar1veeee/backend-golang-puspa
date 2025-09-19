package middlewares

import (
	"backend-golang/internal/adapters/http/types"
	"backend-golang/internal/constants"
	"backend-golang/internal/errors"
	"backend-golang/internal/helpers"
	"backend-golang/internal/infrastructure/config"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{
				Success: false,
				Message: errors.ErrUnauthorized.Error(),
				Errors:  map[string]string{"errors": "Authorization header is required"},
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{
				Success: false,
				Message: errors.ErrUnauthorized.Error(),
				Errors:  map[string]string{"errors": "Authorization header format must be Bearer {token}"},
			})
			return
		}

		tokenString := parts[1]
		claims := &helpers.AppClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWTKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{
				Success: false,
				Message: errors.ErrInvalidToken.Error(),
				Errors:  map[string]string{"errors": "Token is invalid or has expired"},
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), constants.ContextUserID, claims.Subject)
		ctx = context.WithValue(ctx, constants.ContextUserRole, string(claims.Role))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func Authorize(allowedRoles ...constants.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := helpers.GetUserRole(c.Request.Context())
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, types.ErrorResponse{
				Success: false,
				Message: errors.ErrForbidden.Error(),
				Errors:  map[string]string{"errors": "User role not found in token context"},
			})
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == string(role) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, types.ErrorResponse{
				Success: false,
				Message: errors.ErrForbidden.Error(),
				Errors:  map[string]string{"errors": "You are not authorized to access this resource"},
			})
			return
		}

		c.Next()
	}

}
