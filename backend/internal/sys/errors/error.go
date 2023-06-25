package errors

type (
	Err struct {
		msg string
	}
)

var (
	NotImplementedError = NewError("not implemented")
)

func NewError(msg string) Err {
	return Err{msg}
}

func (e Err) Error() string {
	return e.msg
}
