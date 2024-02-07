package exceptions

type NoContentError struct {
	Error string
}

func NewNoContentError(error string) NoContentError {
	return NoContentError{Error: error}
}
