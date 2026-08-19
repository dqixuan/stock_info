package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// THSBoard 同花顺板块
type THSBoard struct {
	BoardType string `json:"板块类型"`
	Code      string `json:"板块代码"`
	Name      string `json:"板块名称"`
}

// THSBoardsResult 对应 Python 返回的 JSON
type THSBoardsResult struct {
	Success       bool       `json:"success"`
	Total         int        `json:"total"`
	IndustryCount int        `json:"industry_count"`
	ConceptCount  int        `json:"concept_count"`
	CSV           string     `json:"csv"`
	Error         string     `json:"error,omitempty"`
	Data          []THSBoard `json:"data"`
}

// GetTHSBoards 调用 Python 脚本获取同花顺板块清单
// boardType: "all" / "industry" / "concept"
func GetTHSBoards(boardType string) (*THSBoardsResult, error) {
	scriptPath, err := resolveScriptPath("get_ths_boards.py")
	if err != nil {
		return nil, err
	}

	args := []string{scriptPath}
	if boardType != "" && boardType != "all" {
		args = append(args, boardType)
	}
	cmd := exec.Command("python3", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Env = stripProxyEnv() // 复用你现有的代理清理函数

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("脚本执行失败: %w\nStderr: %s", err, strings.TrimSpace(stderr.String()))
	}
	if strings.TrimSpace(string(output)) == "" {
		return nil, fmt.Errorf("脚本未返回内容: stderr=%s", strings.TrimSpace(stderr.String()))
	}

	var res THSBoardsResult
	if err := json.Unmarshal(output, &res); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w\n原始输出: %s", err, string(output))
	}
	if !res.Success {
		return nil, fmt.Errorf("脚本返回错误: %s", res.Error)
	}
	return &res, nil
}
