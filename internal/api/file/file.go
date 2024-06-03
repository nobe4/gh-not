package file

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/nobe4/gh-not/internal/gh"
)

type API struct {
	path string
}

func New(path string) *API {
	return &API{path: path}
}

func (a *API) Do(verb string, url string, _ io.Reader, _ interface{}) error {
	return nil
}

func (a *API) Request(verb string, url string, _ io.Reader) (*http.Response, error) {
	if verb == "GET" && url == gh.DefaultEndpoint {
		return a.readFile()
	}

	return nil, errors.New("TODO")
}

func (a *API) readFile() (*http.Response, error) {
	f, err := os.Open(a.path)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Body: io.NopCloser(bufio.NewReader(f)),
	}, nil
}
