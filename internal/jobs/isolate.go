package jobs

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strconv"
)

const isolateVarDir string = "/var/local/lib/isolate"
const isolateBinaryPath string = "/usr/local/bin/isolate"
const cgroupsFlag string = "--cg"
const stderrToStderrFlag string = "--stderr-to-stdout"

func writeStringToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("failed to create source code file")
		return err
	}
	_, err = file.WriteString(content)
	if err != nil {
		slog.Error("failed to write to source code file")
		return err
	}

	return nil

}

func getFileContent(filePath string, n int) string {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("failed to open stdout file", "err", err)
		return ""
	}
	reader := io.Reader(file)
	buf := make([]byte, n)
	_, err = reader.Read(buf)
	if err != nil {
		slog.Error("failed to read stdout file", "err", err)
		return ""
	}

	return string(buf)
}

func initialize(boxId int) error {
	cmd := exec.Command(isolateBinaryPath, "--init", cgroupsFlag, "-b", strconv.Itoa(boxId))
	if err := cmd.Run(); err != nil {
		slog.Error("failed in initialize method", "err", err)
		return err
	}
	slog.Info("initialized isolate box", "boxId", boxId)
	return nil
}

func cleanup(boxId int) error {
	cmd := exec.Command(isolateBinaryPath, "--cleanup", cgroupsFlag, "-b", strconv.Itoa(boxId))
	if err := cmd.Run(); err != nil {
		slog.Error("failed in cleanup method", "err", err)
		return err
	}
	slog.Info("cleaned isolate box", "boxId", boxId)
	return nil
}

func run(boxId int, processCountMax int, command string) {
	const stdoutFileName string = "stdout.txt"
	const stderrFileName string = "stderr.txt"
	cmd := exec.Command(isolateBinaryPath, "--run", cgroupsFlag, "-b", strconv.Itoa(boxId), fmt.Sprintf("-p%d", processCountMax), "-o", stdoutFileName, "-r", stderrFileName, "--", "/bin/sh", "-c", command)
	if err := cmd.Run(); err != nil {
		slog.Error("failed in run method", "err", err)
	}
	slog.Info("run in isolate box", "boxId", boxId)
}

func RunCode() {
	const boxId int = 1
	var workDir string = path.Join(isolateVarDir, strconv.Itoa(boxId))
	var boxDir string = path.Join(workDir, "box")
	var stdoutFileName string = path.Join(boxDir, "stdout.txt")
	var sourceCodeFilePath string = path.Join(boxDir, "main.py")
	var stdoutFileContent string

	initialize(boxId)
	var sourceCodeContent string = "print('Hello From Python')"
	sourceCodeFile, err := os.Create(sourceCodeFilePath)
	if err != nil {
		slog.Error("failed to create source code file")
		return
	}
	_, err = sourceCodeFile.WriteString(sourceCodeContent)
	if err != nil {
		slog.Error("failed to write to source code file")
		return
	}

	run(boxId, 2, "python main.py")
	stdoutFileContent = getFileContent(stdoutFileName, 20)
	slog.Info("Run Result", "stdout", stdoutFileContent)

	cleanup(boxId)

}
