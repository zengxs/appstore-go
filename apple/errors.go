package apple

import (
	"errors"
)

var (
	ErrAlreadyLoggedIn = errors.New("you are already logged in")
	ErrNotLoggedIn     = errors.New("you are not logged in")
	ErrHTTPError       = errors.New("http error occurred")
	ErrNoResults       = errors.New("no results found")
)
