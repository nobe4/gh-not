package github

import (
	gh "github.com/cli/go-gh/v2/pkg/api"
	"github.com/nobe4/gh-not/internal/api"
)

func New() (api.Requestor, error) {
	return gh.DefaultRESTClient()
}
