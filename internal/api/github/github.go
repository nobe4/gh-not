package github

import (
	"fmt"

	gh "github.com/cli/go-gh/v2/pkg/api"

	"github.com/nobe4/gh-not/internal/api"
)

func New() (api.Requestor, error) {
	client, err := gh.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}
