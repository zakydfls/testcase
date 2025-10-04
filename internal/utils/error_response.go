package utils

import "net/http"

type ErrorCode struct {
	Code       int
	Key        string
	Message    string
	HttpStatus int
}

var (
	ErrDataMissing        = ErrorCode{Code: 11, Key: "missing_data", Message: "Data is missing", HttpStatus: http.StatusBadRequest}
	ErrPlatinumField      = ErrorCode{Code: 22, Key: "platinum_field_error", Message: "Platinum field error", HttpStatus: http.StatusBadRequest}
	ErrUnauthorized       = ErrorCode{Code: 33, Key: "unauthorized", Message: "Unauthorized access", HttpStatus: http.StatusUnauthorized}
	ErrInvalidRequest     = ErrorCode{Code: 44, Key: "invalid_request", Message: "Invalid request parameters", HttpStatus: http.StatusBadRequest}
	ErrNotFound           = ErrorCode{Code: 55, Key: "not_found", Message: "Resource not found", HttpStatus: http.StatusNotFound}
	ErrRouteNotFound      = ErrorCode{Code: 66, Key: "route_not_found", Message: "Route not found", HttpStatus: http.StatusNotFound}
	ErrUserListError      = ErrorCode{Code: 77, Key: "user_list_error", Message: "Gagal mengambil data user", HttpStatus: http.StatusInternalServerError}
	ErrUserNotFound       = ErrorCode{Code: 88, Key: "user_not_found", Message: "User tidak ditemukan", HttpStatus: http.StatusNotFound}
	ErrCreateUserError    = ErrorCode{Code: 99, Key: "create_user_error", Message: "Gagal membuat user", HttpStatus: http.StatusInternalServerError}
	ErrFetchDataError     = ErrorCode{Code: 100, Key: "fetch_data_error", Message: "Gagal mengambil data", HttpStatus: http.StatusInternalServerError}
	ErrInvalidCredentials = ErrorCode{Code: 101, Key: "invalid_credentials", Message: "Email atau password salah", HttpStatus: http.StatusUnauthorized}
	ErrTokenExpired       = ErrorCode{Code: 102, Key: "token_expired", Message: "Token telah kadaluarsa", HttpStatus: http.StatusUnauthorized}
	ErrInvalidToken       = ErrorCode{Code: 103, Key: "invalid_token", Message: "Token tidak valid", HttpStatus: http.StatusUnauthorized}
	ErrUpdateDataError    = ErrorCode{Code: 104, Key: "update_data_error", Message: "Gagal memperbarui data", HttpStatus: http.StatusInternalServerError}
	ErrConflict           = ErrorCode{Code: 105, Key: "conflict", Message: "Conflict occurred", HttpStatus: http.StatusConflict}
	ErrInternalServer     = ErrorCode{Code: 106, Key: "internal_server_error", Message: "Internal server error", HttpStatus: http.StatusInternalServerError}
	ErrUsernameExists     = ErrorCode{Code: 107, Key: "username_exists", Message: "Username already exists", HttpStatus: http.StatusConflict}
	ErrEmailExists        = ErrorCode{Code: 108, Key: "email_exists", Message: "Email already exists", HttpStatus: http.StatusConflict}
	ErrInactiveUser       = ErrorCode{Code: 109, Key: "inactive_user", Message: "User is inactive", HttpStatus: http.StatusForbidden}
	ErrForbiddenAccess    = ErrorCode{Code: 110, Key: "forbidden_access", Message: "You do not have permission to access this resource", HttpStatus: http.StatusForbidden}
)

var errorMap = make(map[int]ErrorCode)

func init() {
	registerError(ErrDataMissing)
	registerError(ErrPlatinumField)
	registerError(ErrUnauthorized)
	registerError(ErrInvalidRequest)
	registerError(ErrNotFound)
	registerError(ErrRouteNotFound)
	registerError(ErrUserListError)
	registerError(ErrUserNotFound)
	registerError(ErrCreateUserError)
	registerError(ErrFetchDataError)
	registerError(ErrInvalidCredentials)
	registerError(ErrTokenExpired)
	registerError(ErrInvalidToken)
	registerError(ErrUpdateDataError)
	registerError(ErrConflict)
	registerError(ErrInternalServer)
	registerError(ErrUsernameExists)
	registerError(ErrEmailExists)
	registerError(ErrInactiveUser)
	registerError(ErrForbiddenAccess)
}

func registerError(err ErrorCode) {
	errorMap[err.Code] = err
}

func GetError(err ErrorCode) ErrorCode {
	if err, exists := errorMap[err.Code]; exists {
		return err
	}
	return ErrorCode{Code: err.Code, Key: "unknown_error", Message: "Unknown error"}
}

func GetErrorKey(code int) string {
	return GetError(ErrorCode{Code: code}).Key
}

func GetErrorMessage(code int) string {
	return GetError(ErrorCode{Code: code}).Message
}
