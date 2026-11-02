package ui

type ReportedError struct {
	Err error
}

func (e *ReportedError) Error() string {
	return e.Err.Error()
}

func (e *ReportedError) Unwrap() error {
	return e.Err
}
