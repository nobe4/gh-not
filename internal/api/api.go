package api

import (
	"io"
	"net/http"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

type Caller interface {
	Do(string, string, io.Reader, interface{}) error
	Request(string, string, io.Reader) (*http.Response, error)
}

func NewGH() (Caller, error) {
	return ghapi.DefaultRESTClient()
}
