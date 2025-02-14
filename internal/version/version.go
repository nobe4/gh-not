package version

import (
	"fmt"
	"runtime/debug"
)

const template = "%s (%s) built at %s\nhttps://github.com/nobe4/gh-not/releases/tag/%s"

// Leaving them global to allow setting them via ldflags.
// E.g. go build ./cmd/gh-not -ldflags "-X github.com/nobe4/gh-not/internal/version.tag=0.1.0".
//
//nolint:gochecknoglobals // see above.
var (
	tag    = "UNSET_TAG"
	commit = "UNSET_COMMIT"
	date   = "UNSET_DATE"
)

func String() string {
	if tag == "UNSET_TAG" && commit == "UNSET_COMMIT" && date == "UNSET_DATE" {
		parseBuildInfo()
	}

	return fmt.Sprintf(template, tag, commit, date, tag)
}

func parseBuildInfo() {
	info, ok := debug.ReadBuildInfo()

	if ok {
		tag = info.Main.Version

		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}

			if setting.Key == "vcs.time" {
				date = setting.Value
			}
		}
	}
}
