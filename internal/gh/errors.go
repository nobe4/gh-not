package gh

import (
	"errors"
	"fmt"
)

// decodeError abstracts all decoding error.
var decodeError error = errors.New("decode error")

type RetryError struct {
	verb string
	url  string
}

func (e RetryError) Error() string {
	return fmt.Sprintf("retry exceeded for %s %s", e.verb, e.url)
}
