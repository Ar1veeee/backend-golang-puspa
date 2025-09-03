package errors

import "errors"

var (
	ErrGenerateToken        = errors.New("failed to generate token")
	ErrGenerateRefreshToken = errors.New("failed to generate refresh token")
	ErrSaveRefreshToken     = errors.New("failed to save refresh token")
	ErrTokenExpired         = errors.New("authentication token expired")
	ErrInvalidRefreshToken  = errors.New("invalid or revoked refresh token")
	ErrTooManyLoginAttempts = errors.New("too many loginAttempts")
)
