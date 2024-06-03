package file

import (
	"bytes"
	"errors"
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

func (a *API) Do(verb string, url string, _ io.Reader, _ interface{}) error {
    return nil
}

func (a *API) Request(verb string, url string, _ io.Reader) (*http.Response, error) {
	if verb == "GET" && url == "https://api.github.com/notifications" {
		return a.readFile()
	}

	return nil, errors.New("TODO")
}

func (a *API) readFile() (*http.Response, error) {
	content, err := os.ReadFile(a.path)
    if err != nil {
        return nil, err
    }

    return &http.Response{
        Body: io.NopCloser(bytes.NewBuffer(content)),
    }, nil
}
