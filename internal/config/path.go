package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// inspired by https://github.com/cli/go-gh/blob/trunk/pkg/config/config.go
const (
	appName        = "gh-not"
	appData        = "AppData"
	ghNotConfigDir = "GHNOT_CONFIG_DIR"
	localAppData   = "LocalAppData"
	xdgConfigHome  = "XDG_CONFIG_HOME"
	xdgStateHome   = "XDG_STATE_HOME"
)

// Dir returns the directory where the configuration files are stored.
func Dir() string {
	var path string
	if a := os.Getenv(ghNotConfigDir); a != "" {
		path = a
	} else if b := os.Getenv(xdgConfigHome); b != "" {
		path = filepath.Join(b, appName)
	} else if c := os.Getenv(appData); runtime.GOOS == "windows" && c != "" {
		path = filepath.Join(c, appName)
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", appName)
	}

	return path
}

// StateDir returns the directory where the state files are stored.
func StateDir() string {
	var path string
	if a := os.Getenv(xdgStateHome); a != "" {
		path = filepath.Join(a, appName)
	} else if b := os.Getenv(localAppData); runtime.GOOS == "windows" && b != "" {
		path = filepath.Join(b, appName)
	} else {
		c, _ := os.UserHomeDir()
		path = filepath.Join(c, ".local", "state", appName)
	}

	return path
}
