package errors

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUsernameExists         = errors.New("username already exists")
	ErrEmailExists            = errors.New("email already exists")
	ErrUserIDRequired         = errors.New("user id is required")
	ErrInvalidUserID          = errors.New("invalid user id format")
	ErrInvalidCredentials     = errors.New("invalid user credentials")
	ErrExpiredCode            = errors.New("verification code has expired")
	ErrUserInactive           = errors.New("user account is inactive")
	ErrUserSuspended          = errors.New("user account is suspended")
	ErrInsufficientPermission = errors.New("insufficient permission to perform this action")
	ErrCannotDeleteSelf       = errors.New("cannot delete your own account")
)

var (
	ErrWeakPassword    = errors.New("password must be at least 8 characters")
	ErrPasswordNumber  = errors.New("password must contain at least one number")
	ErrPasswordUpper   = errors.New("password must contain at least one uppercase character")
	ErrPasswordSpecial = errors.New("password must contain at least one special character")
	ErrSamePassword    = errors.New("new password must be different from current password")
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInvalidInput   = errors.New("invalid input provided")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("access forbidden")
	ErrInvalidToken   = errors.New("invalid token")
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource conflict")
	ErrInternalServer = errors.New("internal server errors")
)

var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrTransactionFailed  = errors.New("database transaction failed")
	ErrUniqueViolation    = errors.New("unique constraint violation")
)

var (
	ErrConfigValidation = errors.New("configuration validation failed")
	ErrMissingEnvVar    = errors.New("missing required environment variable")
)
