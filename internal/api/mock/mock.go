package mock

import (
	"fmt"
	"io"
	"net/http"

	"github.com/nobe4/gh-not/internal/api"
)

type Mock struct {
	responses []Response
	calls     []Call

	index int
}

type Response struct {
	Error    error
	Response *http.Response
}

type Call struct {
	verb     string
	endpoint string
}

type MockError struct {
	verb     string
	endpoint string
	message  string
}

func (e *MockError) Error() string {
	return fmt.Sprintf("mock error: %s %s: %s", e.verb, e.endpoint, e.message)
}

func New(r []Response) (api.Caller, error) {
	return &Mock{
		responses: r,
		calls:     []Call{},
		index:     0,
	}, nil
}

func (m *Mock) Do(verb, endpoint string, body io.Reader, out interface{}) error {
	m.calls = append(m.calls, Call{verb, endpoint})

	if m.index >= len(m.responses) {
		return &MockError{verb, endpoint, "no more responses"}
	}

	e := m.responses[m.index].Error
	m.index++

	return e
}

func (m *Mock) Request(verb, endpoint string, body io.Reader) (*http.Response, error) {
	m.calls = append(m.calls, Call{verb, endpoint})

	if m.index >= len(m.responses) {
		return nil, &MockError{verb, endpoint, "no more responses"}
	}

	r := m.responses[m.index]
	m.index++

	return r.Response, r.Error
}
