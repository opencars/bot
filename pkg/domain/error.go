package domain

var (
	// ErrRemoteFailed
	ErrRemoteFailed = NewError("remote failed")
)

type Error struct {
	text string
}

func NewError(text string) Error {
	return Error{
		text: text,
	}
}

func (e Error) Error() string {
	return e.text
}
