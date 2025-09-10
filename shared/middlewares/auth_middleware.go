package middlewares

import (
	"backend-golang/shared/config"
	"backend-golang/shared/constants"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/types"
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
				Message: globalErrors.ErrUnauthorized.Error(),
				Errors:  map[string]string{"errors": "Authorization header is required"},
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{
				Success: false,
				Message: globalErrors.ErrUnauthorized.Error(),
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
				Message: globalErrors.ErrInvalidToken.Error(),
				Errors:  map[string]string{"errors": "Token is invalid or has expired"},
			})
			return
		}

		c.Set("userId", claims.Subject)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func Authorize(allowedRoles ...constants.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, types.ErrorResponse{
				Success: false,
				Message: globalErrors.ErrForbidden.Error(),
				Errors:  map[string]string{"errors": "User role not found in token context"},
			})
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, types.ErrorResponse{
				Success: false,
				Message: globalErrors.ErrForbidden.Error(),
				Errors:  map[string]string{"errors": "You are not authorized to access this resource"},
			})
			return
		}

		c.Next()
	}

}
