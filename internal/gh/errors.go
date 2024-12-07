package gh

import (
	"errors"
	"fmt"
)

// errDecode abstracts all decoding error.
var errDecode = errors.New("decode error")

type RetryError struct {
	verb string
	url  string
}

func (e RetryError) Error() string {
	return fmt.Sprintf("retry exceeded for %s %s", e.verb, e.url)
}
