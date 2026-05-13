package mock

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Mock struct {
	Calls []Call
}

type Error struct {
	verb     string
	endpoint string
	message  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("mock error for call [%s %s]: %s", e.verb, e.endpoint, e.message)
}

func (m *Mock) Done() error {
	for _, c := range m.Calls {
		if !c.Matched {
			return &Error{c.Verb, c.URL, "call not matched"}
		}
	}

	return nil
}

func (m *Mock) Request(verb, endpoint string, _ io.Reader) (*http.Response, error) {
	call, err := m.nextCall(verb, endpoint)
	if err != nil {
		return nil, err
	}

	return call.Response, call.Error
}

func (m *Mock) nextCall(verb, endpoint string) (*Call, error) {
	for i := range m.Calls {
		c := &m.Calls[i]
		if c.Matched {
			continue
		}

		if !c.matches(verb, endpoint) {
			continue
		}

		c.Matched = true

		slog.Debug("mock call", "verb", verb, "endpoint", endpoint, "call", c)

		return c, nil
	}

	return nil, &Error{verb, endpoint, "unexpected call: no match found"}
}
