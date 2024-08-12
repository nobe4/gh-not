package mock

import (
	"fmt"
	"io"
	"net/http"

	"github.com/nobe4/gh-not/internal/api"
)

type Mock struct {
	Calls []Call

	index int
}

type Call struct {
	Verb     string
	Endpoint string
	Data     any
	Error    error
	Response *http.Response
}

type MockError struct {
	verb     string
	endpoint string
	message  string
}

func (e *MockError) Error() string {
	return fmt.Sprintf("mock error: %s %s: %s", e.verb, e.endpoint, e.message)
}

func New(c []Call) api.Requestor {
	return &Mock{Calls: c}
}

func (m *Mock) call(verb, endpoint string) (Call, error) {
	if m.index >= len(m.Calls) {
		return Call{}, &MockError{verb, endpoint, "no more calls"}
	}

	call := m.Calls[m.index]
	if (call.Verb != "" && call.Verb != verb) || (call.Endpoint != "" && call.Endpoint != endpoint) {
		return Call{}, &MockError{verb, endpoint, "unexpected call"}
	}

	m.index++

	return call, nil
}

func (m *Mock) Request(verb, endpoint string, body io.Reader) (*http.Response, error) {
	call, err := m.call(verb, endpoint)
	if err != nil {
		return nil, err
	}

	return call.Response, call.Error
}
