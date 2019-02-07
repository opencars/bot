package openalpr

// APIError is representation of OpenALPR API Error.
type APIError struct {
	msg string
}

// Error implements "error" interface.
func (err *APIError) Error() string {
	return err.msg
}

// New creates new object of "APIError".
func New(text string) error {
	return &APIError{msg: text}
}
