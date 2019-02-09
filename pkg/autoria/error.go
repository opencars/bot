package autoria

// APIError is representation of AUTO.RIA API Error.
type APIError struct {
	msg string
}

// Error implements "error" interface.
func (err *APIError) Error() string {
	return err.msg
}

// NewErr creates new object of "APIError".
func NewErr(text string) error {
	return &APIError{msg: text}
}
