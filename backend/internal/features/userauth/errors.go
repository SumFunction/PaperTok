package userauth

import "errors"

// Common errors for user authentication operations.
var (
	// ErrUserAlreadyExists is returned when attempting to register
	// with an email or username that's already taken.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidCredentials is returned when login credentials are incorrect.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidEmail is returned when the email format is invalid.
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrInvalidUsername is returned when the username format is invalid.
	ErrInvalidUsername = errors.New("invalid username format")

	// ErrWeakPassword is returned when the password doesn't meet security requirements.
	ErrWeakPassword = errors.New("password is too weak")

	// ErrUnauthorized is returned when a user is not authorized to perform an action.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrValidationFailed is returned when input validation fails.
	ErrValidationFailed = errors.New("validation failed")
)

// ErrorCode maps error types to error codes for API responses.
var ErrorCodes = map[error]string{
	ErrUserAlreadyExists:  "USER_EXISTS",
	ErrInvalidCredentials: "INVALID_CREDENTIALS",
	ErrUserNotFound:       "USER_NOT_FOUND",
	ErrInvalidEmail:       "INVALID_EMAIL",
	ErrInvalidUsername:    "INVALID_USERNAME",
	ErrWeakPassword:       "WEAK_PASSWORD",
	ErrUnauthorized:       "UNAUTHORIZED",
	ErrValidationFailed:   "VALIDATION_FAILED",
}

// GetErrorCode returns the error code for a given error.
func GetErrorCode(err error) string {
	if code, ok := ErrorCodes[err]; ok {
		return code
	}
	return "INTERNAL_ERROR"
}

// GetErrorMessage returns a user-friendly error message.
func GetErrorMessage(err error) string {
	switch err {
	case ErrUserAlreadyExists:
		return "该邮箱或用户名已被注册"
	case ErrInvalidCredentials:
		return "邮箱或密码错误"
	case ErrUserNotFound:
		return "用户不存在"
	case ErrInvalidEmail:
		return "邮箱格式不正确"
	case ErrInvalidUsername:
		return "用户名需要3-50个字符，只能包含字母、数字和下划线"
	case ErrWeakPassword:
		return "密码至少8个字符，需包含大写字母、小写字母、数字、特殊字符中的至少2种"
	case ErrUnauthorized:
		return "请先登录"
	case ErrValidationFailed:
		return "输入信息验证失败"
	default:
		return "服务器错误，请稍后重试"
	}
}
