package jobs

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	isolateVarDir     string
	isolateBinaryPath string
)

func getPathOfBinary(binaryName string) (string, error) {

	binaryFilePath, err := exec.LookPath(binaryName)
	if err != nil {
		return binaryFilePath, fmt.Errorf("Could not find binary %s", binaryName)
	}

	slog.Debug("Path to Binary File", "path", binaryFilePath)
	binaryFilePathSymlinkResolved, err := filepath.EvalSymlinks(binaryFilePath)
	if err != nil {
		return binaryFilePathSymlinkResolved, fmt.Errorf("Could not resolve symlink %s", binaryFilePath)
	}

	slog.Debug("Path to Binary File after symlink evaluation", "path", binaryFilePathSymlinkResolved)

	return binaryFilePathSymlinkResolved, nil
}

// exists returns whether the given file or directory exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func InitializeIsolate() error {

	if exists, _ := pathExists("/usr/bin/isolate"); exists {
		isolateBinaryPath = "/usr/bin/isolate"
		isolateVarDir = "/var/lib/isolate"
		slog.Info("isolate binary path found", "path", isolateBinaryPath)
		return nil
	} else if exists, _ := pathExists("/usr/local/bin/isolate"); exists {
		isolateBinaryPath = "/usr/local/bin/isolate"
		isolateVarDir = "/var/local/lib/isolate"
		slog.Info("isolate binary path found", "path", isolateBinaryPath)
		return nil
	} else {
		slog.Error("failed to locate isolate binary")
		return fmt.Errorf("failed to locate isolate binary")
	}
}
