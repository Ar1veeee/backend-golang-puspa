package errors

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUsernameExists         = errors.New("username already exists")
	ErrEmailExists            = errors.New("email already exists")
	ErrUserIDRequired         = errors.New("user id is required")
	ErrInvalidUserID          = errors.New("invalid user id format")
	ErrInvalidCredentials     = errors.New("invalid user credentials")
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
	ErrUserCreationFailed  = errors.New("failed to create user")
	ErrUserUpdateFailed    = errors.New("failed to update user")
	ErrUserDeletionFailed  = errors.New("failed to delete user")
	ErrUserRetrievalFailed = errors.New("failed to retrieve user")
)
