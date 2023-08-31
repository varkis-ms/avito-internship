package apperror

var (
	ErrNoSegment          = New(nil, "The specified segments do not exist or have already been deleted")
	ErrNoUser             = New(nil, "the specified user does not exist")
	ErrBadRequest         = New(nil, "the request to the server contains a syntax error")
	ErrWrongPercent       = New(nil, "percentage must be set in the range 0.0-1.0")
	ErrWrongTtl           = New(nil, "ttl must be strictly positive")
	ErrFileNotFound       = New(nil, "file not found")
	ErrGDriveNotAvailable = New(nil, "Google Drive is unavailable, please try again later")
)

type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"message,omitempty"`
}

func (e *AppError) Error() string { return e.Message }

func New(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
	}
}

func SystemError(err error) *AppError {
	return New(err, "internal system error")
}
