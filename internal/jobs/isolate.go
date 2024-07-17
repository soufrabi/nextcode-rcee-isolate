package jobs

import (
	"fmt"
	"git.soufrabi.com/nextcode/rcee-isolate/internal/api"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const cgroupsFlag string = "--cg"
const stderrToStderrFlag string = "--stderr-to-stdout"

func getMetadataMap(filePath string) map[string]string {
	var m map[string]string
	m = make(map[string]string)
	var metadataFileContent string = getFileContent(filePath, 128)
	var lines []string = strings.Split(metadataFileContent, "\n")
	slog.Debug("Metadata Lines", "lines", lines)
	for _, line := range lines {
		slog.Debug("Metadata Line", "line", line)
		var keyVal []string = strings.Split(line, ":")
		if len(keyVal) == 2 {
			var key string = keyVal[0]
			var val string = keyVal[1]
			if len(key) > 0 && len(val) > 0 {
				m[key] = val
			}
		}
	}

	return m
}

func writeStringToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("failed to create file", "path", filePath, "err", err)
		return err
	}
	_, err = file.WriteString(content)
	if err != nil {
		slog.Error("failed to write to file", "path", filePath, "err", err)
		return err
	}

	return nil

}

func convertMillisecondsToSecondsInString(n uint64) string {
	return fmt.Sprintf("%d.%03d", n/1_000, n%1_000)
}

func getFileContent(filePath string, maxBytes uint64) string {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("failed to open file", "path", filePath, "err", err)
		return ""
	}
	reader := io.Reader(file)
	buf := make([]byte, maxBytes)
	n, err := reader.Read(buf)
	if err != nil {
		slog.Error("failed to read file", "path", filePath, "err", err)
		return ""
	}

	return string(buf[:n])
}

func initialize(boxId int) error {
	cmd := exec.Command(isolateBinaryPath, cgroupsFlag, "--init", "-b", strconv.Itoa(boxId))
	if err := cmd.Run(); err != nil {
		slog.Error("failed in initialize method", "err", err)
		return err
	}
	slog.Info("initialized isolate box", "boxId", boxId)
	return nil
}

func cleanup(boxId int) error {
	cmd := exec.Command(isolateBinaryPath, cgroupsFlag, "--cleanup", "-b", strconv.Itoa(boxId))
	if err := cmd.Run(); err != nil {
		slog.Error("failed in cleanup method", "err", err)
		return err
	}
	slog.Info("cleaned isolate box", "boxId", boxId)
	return nil
}

func run(boxId int, req api.RunRequest, command string, stdoutFileName string, stderrFileName string, metadataFileName string) error {
	cmd := exec.Command(
		isolateBinaryPath,
		cgroupsFlag,
		"-s",
		"-b",
		fmt.Sprintf("%v", boxId),
		"-M",
		metadataFileName,
		"-t",
		convertMillisecondsToSecondsInString(req.CpuTimeLimit),
		"-x",
		convertMillisecondsToSecondsInString(req.CpuExtraTime),
		"-w",
		convertMillisecondsToSecondsInString(req.WallTimeLimit),
		fmt.Sprintf("-p%d", req.MaxProcessesAndOrThreads),
		fmt.Sprintf("--cg-mem=%v", req.MemoryLimit),
		"-k",
		fmt.Sprintf("%v", req.StackLimit),
		"-f",
		fmt.Sprintf("%v", req.MaxFileSize),
		"-o",
		stdoutFileName,
		"-r",
		stderrFileName,
		"--run",
		"--",
		"/bin/sh",
		"-c",
		command,
	)
	if err := cmd.Run(); err != nil {
		slog.Error("failed in run method", "err", err)
		return err
	}
	slog.Info("run in isolate box", "boxId", boxId)
	return nil
}

func generateErrorResponse(message string) api.RunResponse {
	return api.RunResponse{
		Stdout:   "",
		Stderr:   "",
		Status:   message,
		ExitCode: "-1",
		Time:     "0",
		WallTime: "0",
		Memory:   "0",
	}

}

func RunCode(request api.RunRequest) api.RunResponse {
	var err error
	var metadataTmpFile *os.File
	metadataTmpFile, err = os.CreateTemp("", "metadata*.txt")
	defer os.Remove(metadataTmpFile.Name())
	slog.Debug("Metadata Temp File", "path", metadataTmpFile.Name())

	const boxId int = 1
	var (
		coreDir            string = path.Join(isolateVarDir, strconv.Itoa(boxId))
		boxDir             string = path.Join(coreDir, "box")
		stdoutFileName     string = "stdout.txt"
		stderrFileName     string = "stderr.txt"
		metadataFileName   string = metadataTmpFile.Name()
		stdoutFilePath     string = path.Join(boxDir, stdoutFileName)
		stderrFilePath     string = path.Join(boxDir, stderrFileName)
		metadataFilePath   string = metadataFileName
		sourceCodeFilePath string = path.Join(boxDir, "main.py")
		stdoutFileContent  string
		stderrFileContent  string
		metadataMap        map[string]string
	)

	err = initialize(boxId)
	defer cleanup(boxId)
	if err != nil {
		return generateErrorResponse("INTERNAL ERROR")
	}

	err = writeStringToFile(sourceCodeFilePath, request.SourceCode)
	if err != nil {
		return generateErrorResponse("INTERNAL ERROR")
	}

	run(boxId, request, "python main.py", stdoutFileName, stderrFileName, metadataFileName)
	stdoutFileContent = getFileContent(stdoutFilePath, request.MaxFileSize)
	stderrFileContent = getFileContent(stderrFilePath, request.MaxFileSize)
	metadataMap = getMetadataMap(metadataFilePath)
	slog.Info("Metadata Map", "content", metadataMap)

	res := api.RunResponse{
		Stdout: stdoutFileContent,
		Stderr: stderrFileContent,
		Status: "",
	}
	slog.Debug("Run Result", "res", res)
	return res

}
