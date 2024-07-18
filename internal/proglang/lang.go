package proglang

type Language struct {
	Id                uint64
	CompilationNeeded bool
	InterpreterNeeded bool
	FileExtension     string
	InterpreterName   string
}

const (
	PYTHON_CPYTHON_3_12  uint64 = 1
	RUBY_3_2             uint64 = 2
	JAVASCRIPT_NODEJS_18 uint64 = 3
)

var LangMap = map[uint64]Language{
	PYTHON_CPYTHON_3_12: {
		CompilationNeeded: false,
		InterpreterNeeded: true,
		FileExtension:     "py",
		InterpreterName:   "python3.12",
	},
	RUBY_3_2: {
		CompilationNeeded: false,
		InterpreterNeeded: true,
		FileExtension:     "rb",
		InterpreterName:   "ruby3.2",
	},
	JAVASCRIPT_NODEJS_18: {
		CompilationNeeded: false,
		InterpreterNeeded: true,
		FileExtension:     "js",
		InterpreterName:   "node",
	},
}
