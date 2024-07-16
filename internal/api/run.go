package api

type RunRequest struct {
	LanguageId               uint64 `json:"languageId"`
	SourceCode               string `json:"sourceCode"`
	Stdin                    string `json:"stdin"`
	CpuTimeLimit             uint64 `json:"cpuTimeLimit"`
	CpuExtraTime             uint64 `json:"cpuExtraTime"`
	WallTimeLimit            uint64 `json:"wallTimeLimit"`
	MaxProcessesAndOrThreads uint64 `json:"maxProcessesAndOrThreads"`
	MemoryLimit              uint64 `json:"memoryLimit"`
	StackLimit               uint64 `json:"stackLimit"`
	MaxFileSize              uint64 `json:"maxFileSize"`
}

type RunResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Status string `json:"status"`
}
