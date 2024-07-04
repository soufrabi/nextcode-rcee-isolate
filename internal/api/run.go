package api

type RunRequest struct {
    SourceCode string `json:"sourceCode"`
    Stdin      string `json:"stdin"`
    CpuTimeLimit uint16 `json:"cpuTimeLimit"`
    MemoryLimit uint16 `json:"memoryLimit"`
}

type RunResponse struct {
    Stdout string `json:"stdout"`
    Stderr string `json:"stderr"`
    Status string `json:"status"`
}

