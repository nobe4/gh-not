package version

import (
	"fmt"
	"runtime/debug"
)

func String() string {
	const template = "%s (%s) built at %s\nhttps://github.com/nobe4/gh-not/releases/tag/%s"

	var (
		tag    = "UNSET_TAG"
		commit = "UNSET_COMMIT"
		date   = "UNSET_DATE"
	)

	info, ok := debug.ReadBuildInfo()

	if ok {
		fmt.Printf("%+v\n", info)

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

	return fmt.Sprintf(template, tag, commit, date, tag)
}
