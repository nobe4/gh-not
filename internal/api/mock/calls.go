package mock

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

type Call struct {
	Verb     string
	Endpoint string
	Data     any
	Error    error
	Response *http.Response
}

type RawCall struct {
	Verb     string      `json:"verb"`
	Endpoint string      `json:"endpoint"`
	Data     any         `json:"data"`
	Error    error       `json:"error"`
	Response RawResponse `json:"response"`
}

type RawResponse struct {
	Headers    map[string][]string `json:"headers"`
	StatusCode int                 `json:"status_code"`
	Body       any                 `json:"body"`
}

func LoadCallsFromFile(path string) ([]Call, error) {
	rawCalls := []RawCall{}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, &rawCalls); err != nil {
		return nil, err
	}

	calls := make([]Call, len(rawCalls))
	for i, call := range rawCalls {
		body, err := json.Marshal(call.Response.Body)
		if err != nil {
			return nil, err
		}

		calls[i] = Call{
			Verb:     call.Verb,
			Endpoint: call.Endpoint,
			Data:     call.Data,
			Error:    call.Error,
			Response: &http.Response{
				Header:     http.Header(call.Response.Headers),
				StatusCode: call.Response.StatusCode,
				Body:       io.NopCloser(strings.NewReader(string(body))),
			},
		}
	}

	return calls, nil
}
