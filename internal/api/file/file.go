package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type API struct {
	path string
}

func New(path string) *API {
	return &API{path: path}
}

func (a *API) Request(verb string, url string, _ io.Reader) (*http.Response, error) {
	if verb == "GET" {
		return a.readFile()
	}

	return nil, errors.New("TODO")
}

func (a *API) readFile() (*http.Response, error) {
	f, err := os.Open(a.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &http.Response{
		Body: io.NopCloser(bufio.NewReader(f)),
	}, nil
}
