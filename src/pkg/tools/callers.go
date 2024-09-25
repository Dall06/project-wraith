package tools

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func ExtractCallerInfo(level int) string {
	// Get caller info, 2 because we want the caller of the function that calls this helper
	pc, file, line, ok := runtime.Caller(level)
	if !ok {
		return "unknown:0 unknown()"
	}

	funcName := runtime.FuncForPC(pc).Name()
	file = filepath.Base(file) // Only use the base name of the file

	// Return formatted string with file, line, and method
	return fmt.Sprintf("%s:%d %s()", file, line, funcName)
}
