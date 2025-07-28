package constant

import "errors"

var (
	ErrLoginFailed                = errors.New("login failed")
	ErrTokenParse                 = errors.New("token parse error")
	ErrMissingOrInvalidAuthHeader = errors.New("authorization header is missing or invalid")
	ErrTokenUsedElsewhere         = errors.New("token expired or used elsewhere")
	ErrUnauthorizedOrNotLoggedIn  = errors.New("user unauthorized or not logged in")
	ErrUUIDIsEmpty                = errors.New("uuid is empty")
)
