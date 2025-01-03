package mock

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Call struct {
	Verb     string
	URL      string
	Data     any
	Error    error
	Response *http.Response
}

type RawCall struct {
	Verb     string      `json:"verb"`
	URL      string      `json:"endpoint"`
	Data     any         `json:"data"`
	Error    RawError    `json:"error"`
	Response RawResponse `json:"response"`
}

type RawResponse struct {
	Headers    map[string][]string `json:"headers"`
	StatusCode int                 `json:"status_code"`
	Body       any                 `json:"body"`
}

type RawError struct {
	StatusCode int `json:"status_code"`
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

	for i, rawCall := range rawCalls {
		body, err := json.Marshal(rawCall.Response.Body)
		if err != nil {
			return nil, err
		}

		call := Call{
			Verb: rawCall.Verb,
			URL:  rawCall.URL,
			Data: rawCall.Data,
			Response: &http.Response{
				Header:     http.Header(rawCall.Response.Headers),
				StatusCode: rawCall.Response.StatusCode,
				Body:       io.NopCloser(strings.NewReader(string(body))),
			},
		}

		if rawCall.Error.StatusCode != 0 {
			call.Error = &api.HTTPError{StatusCode: rawCall.Error.StatusCode}
		}

		calls[i] = call
	}

	return calls, nil
}
