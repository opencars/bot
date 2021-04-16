package domain

var (
	ErrNotRecognized = NewError("number plate is not recognized")
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
