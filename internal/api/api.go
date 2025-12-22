//revive:disable:var-naming // This is fine.
package api

import (
	"io"
	"net/http"
)

type Requestor interface {
	Request(method string, path string, body io.Reader) (*http.Response, error)
}
