package error


type baseError struct {
	err string
}

func (be baseError) Error () string {
	return be.err
}

type NaluError struct {
	baseError
}

func NewNaluError(err string) NaluError {
	return NaluError{baseError{err}}
}

