package version

import (
	"fmt"
	"runtime/debug"
)

var tag = "dev" // set via ldflags
var commit = "123abc"
var date = "now"

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

	return fmt.Sprintf("%s (%s) built at %s\nhttps://github.com/nobe4/gh-not/releases/tag/%s", tag, commit, date, tag)
}
