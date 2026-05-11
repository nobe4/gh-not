package mock

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

type Mock struct {
	Calls []Call

	mu      sync.Mutex
	matched []bool // tracks which calls have been consumed
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
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureMatched()

	for i, call := range m.Calls {
		if !m.matched[i] {
			return &Error{call.Verb, call.URL, fmt.Sprintf("%d calls remaining", m.remaining())}
		}
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

func (m *Mock) remaining() int {
	count := 0

	for _, done := range m.matched {
		if !done {
			count++
		}
	}

	return count
}

func (m *Mock) ensureMatched() {
	if m.matched == nil {
		m.matched = make([]bool, len(m.Calls))
	}
}

func callMatches(call Call, verb, endpoint string) bool {
	if call.Verb != "" && call.Verb != verb {
		return false
	}

	if call.URL != "" && call.URL != endpoint {
		return false
	}

	return true
}

func (m *Mock) nextMatchingCall(verb, endpoint string) int {
	for i, call := range m.Calls {
		if m.matched[i] {
			continue
		}

		if !callMatches(call, verb, endpoint) {
			continue
		}

		return i
	}

	return -1
}

func (m *Mock) call(verb, endpoint string) (Call, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureMatched()

	idx := m.nextMatchingCall(verb, endpoint)
	if idx < 0 {
		return Call{}, &Error{verb, endpoint, "unexpected call: no matching call found"}
	}

	call := m.Calls[idx]
	m.matched[idx] = true

	slog.Debug("mock call", "verb", verb, "endpoint", endpoint, "call", call)

	return call, nil
}
