package mock

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Mock struct {
	Calls []Call

	index int
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
	if m.index < len(m.Calls) {
		return &Error{"", "", fmt.Sprintf("%d calls remaining", len(m.Calls)-m.index)}
	}

	return nil
}

func (m *Mock) Request(verb, endpoint string, _ io.Reader) (*http.Response, error) {
	call, err := m.call(verb, endpoint)
	if err != nil {
		return nil, err
	}

	return call.Response, call.Error
}

func (m *Mock) call(verb, endpoint string) (Call, error) {
	if m.index >= len(m.Calls) {
		return Call{}, &Error{verb, endpoint, "unexpected call: no more calls"}
	}

	call := m.Calls[m.index]
	if (call.Verb != "" && call.Verb != verb) || (call.URL != "" && call.URL != endpoint) {
		return Call{}, &Error{
			verb,
			endpoint,
			fmt.Sprintf("unexpected call: mismatch, expected [%s %s]", call.Verb, call.URL),
		}
	}

	m.index++

	slog.Debug("mock call", "verb", verb, "endpoint", endpoint, "call", call)

	return call, nil
}
