package exceptions

type TokenExpiredException struct {
	Message string
}

func (t *TokenExpiredException) Error() string {
	return t.Message
}

type TokenNotFoundException struct {
	Message string
}

func (t *TokenNotFoundException) Error() string {
	return t.Message
}
