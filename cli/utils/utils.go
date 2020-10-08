package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// ValidateOutputPath validates the output paths of the `export` and `save` commands.
func ValidateOutputPath(path string) error {
	dir := filepath.Dir(filepath.Clean(path))
	if dir != "" && dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return errors.New("invalid output path: directory " + dir + " does not exist")
		}
	}
	// check whether `path` points to a regular file
	// (if the path exists and doesn't point to a directory)
	if fileInfo, err := os.Stat(path); !os.IsNotExist(err) {
		if err != nil {
			return err
		}

		if fileInfo.Mode().IsDir() || fileInfo.Mode().IsRegular() {
			return nil
		}

		if err := ValidateOutputPathFileMode(fileInfo.Mode()); err != nil {
			return errors.New("invalid output path: " + path + " must be a directory or a regular file")
		}
	}
	return nil
}

// ValidateOutputPathFileMode validates the output paths of the `cp` command and serves as a
// helper to `ValidateOutputPath`
func ValidateOutputPathFileMode(fileMode os.FileMode) error {
	switch {
	case fileMode&os.ModeDevice != 0:
		return errors.New("got a device")
	case fileMode&os.ModeIrregular != 0:
		return errors.New("got an irregular file")
	}
	return nil
}

// GetHomeDir gets the home directory. If we're in CI, then it's the CWD, otherwise assume it's
// /home/seluser
func GetHomeDir() (string, error) {
	if (os.Getenv("CI") == "") {
		return "/home/seluser", nil
	}
	workingDir, err := os.Getwd()
	if (err != nil) {
		return "", err
	}
	return workingDir, nil
}

// GetConfigFile config yaml. If we're in CI, then it's config.yaml, otherwise it's config-local.yaml
func GetConfigFile () string {
	if (os.Getenv("CI") == "") {
		return "config.yaml"
	}
	return "config-local.yaml"
}