package api

import (
	"io"
	"net/http"
)

type Requestor interface {
	Request(string, string, io.Reader) (*http.Response, error)
}
