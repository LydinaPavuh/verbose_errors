package error_three

type wError struct {
	msg     string
	wrapped []error
}

func (e *wError) Error() string {
	return e.msg
}
func (e *wError) Unwrap() []error {
	return e.wrapped
}

// Error new error with msg
func Error(msg string, err ...error) error {
	return &wError{
		msg:     msg,
		wrapped: err,
	}
}
