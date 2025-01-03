package api

import (
	"io"
	"net/http"
)

type Requestor interface {
	Request(method string, url string, body io.Reader) (*http.Response, error)
}
