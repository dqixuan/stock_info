package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// StockData 结构体，对应 Python 脚本返回的 JSON 字段
type StockData struct {
	Symbol           string  `json:"symbol"`
	Name             string  `json:"name"`
	LatestPrice      float64 `json:"latest_price"`
	ChangeAmount     float64 `json:"change_amount"`
	ChangePercentage float64 `json:"change_percentage"`
	BuyPrice         float64 `json:"buy_price"`
	SellPrice        float64 `json:"sell_price"`
	LastClose        float64 `json:"last_close"`
	OpenPrice        float64 `json:"open_price"`
	HighPrice        float64 `json:"high_price"`
	LowPrice         float64 `json:"low_price"`
	Volume           float64 `json:"volume"`
	Turnover         float64 `json:"turnover"`
	Timestamp        string  `json:"timestamp"`
	Currency         string  `json:"currency"`
	Error            string  `json:"error,omitempty"`
}

// GetStockPrice 调用 Python 脚本获取单只股票实时行情
func GetStockPrice(symbol string) (*StockData, error) {
	now := time.Now()
	scriptPath, err := resolveScriptPath("get_stock_price.py")
	if err != nil {
		return nil, err
	}

	// 调用 Python 脚本，传入股票代码
	cmd := exec.Command("python3", scriptPath, symbol)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// 关键：清掉代理环境变量，避免 requests 走不可用的代理访问东财
	cmd.Env = stripProxyEnv()

	// 执行命令并获取标准输出
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("脚本执行失败: %w\nStderr: %s", err, strings.TrimSpace(stderr.String()))
	}
	if strings.TrimSpace(string(output)) == "" {
		return nil, fmt.Errorf("脚本未返回内容: symbol=%s stderr=%s", symbol, strings.TrimSpace(stderr.String()))
	}

	// 解析 JSON 输出
	var stockData StockData
	if err := json.Unmarshal(output, &stockData); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w\n原始输出: %s", err, string(output))
	}

	// 检查 Python 脚本是否返回了错误信息
	if stockData.Error != "" {
		return nil, fmt.Errorf("脚本返回错误: %s", stockData.Error)
	}
	fmt.Printf("GetStockPrice for %s cost time: %d(s)\n", symbol, int(time.Since(now).Seconds()))
	return &stockData, nil
}

// stripProxyEnv 返回清理掉所有代理变量后的环境变量列表
func stripProxyEnv() []string {
	var out []string
	for _, kv := range osEnviron() {
		key := kv
		if i := strings.IndexByte(kv, '='); i >= 0 {
			key = kv[:i]
		}
		if strings.Contains(strings.ToLower(key), "proxy") {
			continue // 丢弃 http_proxy / all_proxy / no_proxy 等全部代理变量
		}
		out = append(out, kv)
	}
	return out
}

// osEnviron 封装 os.Environ()，便于测试替换
func osEnviron() []string {
	return osEnvironImpl()
}
