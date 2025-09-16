package errors

import "net/http"

type HTTPError struct {
	Code       string
	StatusCode int
	Message    string
	UserMsg    string
}

func (e HTTPError) Error() string {
	return e.Code
}

func BadRequest(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusBadRequest, "Bad Request", userMsg}
}

func Unauthorized(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusUnauthorized, "Unauthorized", userMsg}
}

func Forbidden(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusForbidden, "Forbidden", userMsg}
}

func NotFound(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusNotFound, "Not Found", userMsg}
}

func Conflict(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusConflict, "Data Conflict", userMsg}
}

func ValidationError(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusUnprocessableEntity, "Validation Errors", userMsg}
}

func InternalServer(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusInternalServerError, "Internal Server Error", userMsg}
}

func TooManyRequests(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusTooManyRequests, "Too Many Requests", userMsg}
}

func Locked(code, userMsg string) HTTPError {
	return HTTPError{code, http.StatusLocked, "Locked", userMsg}
}
