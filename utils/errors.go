package utils

type LastFMError struct {
	Message string
}

func (e LastFMError) Error() string {
	return e.Message
}
