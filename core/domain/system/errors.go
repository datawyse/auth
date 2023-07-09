package system

var NotFoundError = "not found"
var AlreadyExistsError = "already exists"
var InvalidInputError = "invalid input"
var InvalidTokenError = "invalid token"
var ServerError = "server error"
var UserAlreadyExistsError = "user already exists"
var UserNotFoundError = "user not found"
var InvalidCredentialsError = "invalid credentials"
var InvalidUUIDError = "invalid uuid"
var AuthorizationError = "authorization error"

// ApiError - create base error
type ApiError struct {
	Code     int    `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	RawError error  `json:"rawError,omitempty"`
}

// Error - implement error interface
func (e *ApiError) Error() string {
	return e.Message
}

// ErrInvalidUUID
var ErrInvalidUUID = &ApiError{Code: 400, Message: InvalidUUIDError}

// ErrNotFound is the error returned when a resource is not found.
var ErrNotFound = &ApiError{Code: 404, Message: NotFoundError}

// ErrAlreadyExists is the error returned when a resource already exists.
var ErrAlreadyExists = &ApiError{Code: 409, Message: AlreadyExistsError}

// ErrInvalidInput is the error returned when the input is invalid.
var ErrInvalidInput = &ApiError{Code: 400, Message: InvalidInputError}

// ErrInvalidCredentials is the error returned when the credentials are invalid.
var ErrInvalidCredentials = &ApiError{Code: 401, Message: InvalidCredentialsError}

// ErrInvalidToken is the error returned when the token is invalid.
var ErrInvalidToken = &ApiError{Code: 401, Message: InvalidTokenError}

// ErrServer is the error returned when the server has an error.
var ErrServer = &ApiError{Code: 500, Message: ServerError}

var ErrUserAlreadyExists = &ApiError{Code: 409, Message: UserAlreadyExistsError}

var ErrUserNotFound = &ApiError{Code: 404, Message: UserNotFoundError}

var ErrAuthorization = &ApiError{Code: 401, Message: AuthorizationError}

//
// // ErrNotFound is the error returned when a resource is not found.
// var ErrNotFound = errors.New("not found")
//
// // ErrAlreadyExists is the error returned when a resource already exists.
// var ErrAlreadyExists = errors.New("already exists")
//
// // ErrInvalidInput is the error returned when the input is invalid.
// var ErrInvalidInput = errors.New("invalid input")
//
// // ErrUserAlreadyExists is the error returned when a user already exists.
// var ErrUserAlreadyExists = errors.New("user already exists")
//
// // ErrInvalidAuthenticationToken is the error returned when the authentication token is invalid.
// var ErrInvalidAuthenticationToken = errors.New("invalid authentication token")
//
// // ErrInvalidAppContext is the error returned when the app context is invalid.
// var ErrInvalidAppContext = errors.New("invalid app context")
//
// // ErrInvalidCredentials is the error returned when the credentials are invalid.
// var ErrInvalidCredentials = errors.New("invalid credentials")
//
// // ErrInvalidPassword is the error returned when the password is invalid.
// var ErrInvalidPassword = errors.New("invalid password")
//
// // ErrInvalidEmail is the error returned when the email is invalid.
// var ErrInvalidEmail = errors.New("invalid email")
//
// // ErrInvalidUsername is the error returned when the username is invalid.
// var ErrInvalidUsername = errors.New("invalid username")
