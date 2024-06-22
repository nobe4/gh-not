package api

import (
	"io"
	"net/http"
)

type Caller interface {
	Do(string, string, io.Reader, interface{}) error
	Request(string, string, io.Reader) (*http.Response, error)
}
