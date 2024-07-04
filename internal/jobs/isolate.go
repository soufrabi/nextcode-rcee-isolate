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
)

const isolateVarDir string = "/var/local/lib/isolate"
const isolateBinaryPath string = "/usr/local/bin/isolate"
const cgroupsFlag string = "--cg"
const stderrToStderrFlag string = "--stderr-to-stdout"

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

func getFileContent(filePath string, n int) string {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("failed to open file", "path", filePath, "err", err)
		return ""
	}
	reader := io.Reader(file)
	buf := make([]byte, n)
	_, err = reader.Read(buf)
	if err != nil {
		slog.Error("failed to read file", "path", filePath, "err", err)
		return ""
	}

	return string(buf)
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
		fmt.Sprintf("%v", req.CpuTimeLimit),
		"-x",
		fmt.Sprintf("%v", req.CpuExtraTime),
		"-w",
		fmt.Sprintf("%v", req.WallTimeLimit),
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

func RunCode(request api.RunRequest) api.RunResponse {
	var err error
	var metadataTmpFile *os.File
	metadataTmpFile, err = os.CreateTemp("", "metadata*.txt")
	defer os.Remove(metadataTmpFile.Name())
	slog.Debug("Metadata Temp File", "path", metadataTmpFile.Name())

	const boxId int = 1
	var (
		coreDir             string = path.Join(isolateVarDir, strconv.Itoa(boxId))
		boxDir              string = path.Join(coreDir, "box")
		stdoutFileName      string = "stdout.txt"
		stderrFileName      string = "stderr.txt"
		metadataFileName    string = metadataTmpFile.Name()
		stdoutFilePath      string = path.Join(boxDir, stdoutFileName)
		stderrFilePath      string = path.Join(boxDir, stderrFileName)
		metadataFilePath    string = metadataFileName
		sourceCodeFilePath  string = path.Join(boxDir, "main.py")
		stdoutFileContent   string
		stderrFileContent   string
		metadataFileContent string
	)

	err = initialize(boxId)
	defer cleanup(boxId)
	if err != nil {
		return api.RunResponse{
			Stdout: "",
			Stderr: "",
			Status: "INTERNAL_ERROR",
		}
	}

	err = writeStringToFile(sourceCodeFilePath, request.SourceCode)
	if err != nil {
		return api.RunResponse{
			Stdout: "",
			Stderr: "",
			Status: "INTERNAL_ERROR",
		}
	}

	run(boxId, request, "python main.py", stdoutFileName, stderrFileName, metadataFileName)
	stdoutFileContent = getFileContent(stdoutFilePath, 50)
	stderrFileContent = getFileContent(stderrFilePath, 50)
	metadataFileContent = getFileContent(metadataFilePath, 50)
	slog.Info("Metadata file", "content", metadataFileContent)

	res := api.RunResponse{
		Stdout: stdoutFileContent,
		Stderr: stderrFileContent,
		Status: "",
	}
	slog.Debug("Run Result", "res", res)
	return res

}
