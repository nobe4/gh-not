package gh

import "fmt"

type RetryError struct {
	verb     string
	endpoint string
}

func (e RetryError) Error() string {
	return fmt.Sprintf("retry exceeded for %s %s", e.verb, e.endpoint)
}
