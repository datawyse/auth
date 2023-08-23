package system

import "errors"

var NotFoundError = errors.New("not found")
var InvalidInputError = errors.New("invalid input")
var InvalidTokenError = errors.New("invalid token")
var ServerError = errors.New("server error")
var UserAlreadyExistsError = errors.New("user already exists")
var UserNotFoundError = errors.New("user not found")
var InvalidCredentialsError = errors.New("invalid credentials")
var InvalidUUIDError = errors.New("invalid uuid")
var AuthorizationError = errors.New("authorization error")
var InvalidAuthorizationToken = errors.New("invalid authorization token")
var RetrospectTokenError = errors.New("retrospect token error")
var PermissionDeniedError = errors.New("permission denied")
var ValidationError = errors.New("validation error")

// ApiError - create base error
type ApiError struct {
	Code     int    `json:"code,omitempty"`
	Status   bool   `json:"status,omitempty"`
	Message  string `json:"message,omitempty"`
	RawError error  `json:"rawError,omitempty"`
}

// Error - implement error interface
func (e *ApiError) Error() string {
	return e.Message
}

// ErrInvalidUUID is the error returned when a uuid is invalid.
var ErrInvalidUUID = &ApiError{Code: 400, Message: InvalidUUIDError.Error()}

// ErrNotFound is the error returned when a resource is not found.
var ErrNotFound = &ApiError{Code: 404, Message: NotFoundError.Error()}

// ErrInvalidInput is the error returned when the input is invalid.
var ErrInvalidInput = &ApiError{Code: 400, Message: InvalidInputError.Error()}

// ErrInvalidCredentials is the error returned when the credentials are invalid.
var ErrInvalidCredentials = &ApiError{Code: 401, Message: InvalidCredentialsError.Error()}

// ErrInvalidToken is the error returned when the token is invalid.
var ErrInvalidToken = &ApiError{Code: 401, Message: InvalidTokenError.Error()}

// ErrServer is the error returned when the server has an error.
var ErrServer = &ApiError{Code: 500, Message: ServerError.Error()}

var ErrUserAlreadyExists = &ApiError{Code: 409, Message: UserAlreadyExistsError.Error()}

var ErrUserNotFound = &ApiError{Code: 404, Message: UserNotFoundError.Error()}

var ErrAuthorization = &ApiError{Code: 401, Message: AuthorizationError.Error()}

var ErrInvalidAuthorizationToken = &ApiError{Code: 401, Message: InvalidAuthorizationToken.Error()}

var ErrRetrospectToken = &ApiError{Code: 401, Message: RetrospectTokenError.Error()}

var ErrPermissionDenied = &ApiError{Code: 401, Message: PermissionDeniedError.Error()}

var ErrValidationError = &ApiError{Code: 400, Message: ValidationError.Error()}
