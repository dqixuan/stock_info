package pkg

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func resolveScriptPath(scriptName string) (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("获取脚本目录失败")
	}
	return filepath.Join(filepath.Dir(currentFile), "..", "scripts", scriptName), nil
}
