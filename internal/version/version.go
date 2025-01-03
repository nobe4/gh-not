package version

import (
	"fmt"
	"runtime/debug"
)

var (
	tag    = "UNSET_TAG" // set via ldflags
	commit = "UNSET_COMMIT"
	date   = "UNSET_DATE"
)

const template = "%s (%s) built at %s\nhttps://github.com/nobe4/gh-not/releases/tag/%s"

func String() string {
	info, ok := debug.ReadBuildInfo()

	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}

			if setting.Key == "vcs.time" {
				date = setting.Value
			}
		}
	}

	return fmt.Sprintf(template, tag, commit, date, tag)
}
