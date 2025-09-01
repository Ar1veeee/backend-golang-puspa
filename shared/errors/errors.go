package errors

import "errors"

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
