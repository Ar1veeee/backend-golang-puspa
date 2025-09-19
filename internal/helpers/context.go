package helpers

import (
	"backend-golang/internal/constants"
	"context"
)

func GetUserID(ctx context.Context) (string, bool) {
	val := ctx.Value(constants.ContextUserID)
	if val == nil {
		return "", false
	}
	userId, ok := val.(string)
	return userId, ok
}

func GetUserRole(ctx context.Context) (string, bool) {
	val := ctx.Value(constants.ContextUserRole)
	if val == nil {
		return "", false
	}
	
	role, ok := val.(string)
	return role, ok
}
