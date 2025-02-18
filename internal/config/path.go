package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

var errTildeUsage = errors.New("tilde in path is not supported, use $HOME instead")

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

func ExpandPathWithoutTilde(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		return "", fmt.Errorf("%w: config path: %s", errTildeUsage, path)
	}

	// Allows to use $HOME and other environment variables in the configuration
	// paths.
	return os.ExpandEnv(path), nil
}
