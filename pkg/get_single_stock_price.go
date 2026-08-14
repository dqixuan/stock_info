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

func printStockData(stock *StockData) {
	fmt.Println("========================================")
	fmt.Printf("  股票代码: %s\n", stock.Symbol)
	fmt.Printf("  股票名称: %s\n", stock.Name)
	fmt.Println("========================================")
	fmt.Printf("  最新价格: %.2f %s\n", stock.LatestPrice, stock.Currency)
	fmt.Printf("  涨跌额:   %.2f\n", stock.ChangeAmount)
	fmt.Printf("  涨跌幅:   %.3f%%\n", stock.ChangePercentage)
	fmt.Println("----------------------------------------")
	fmt.Printf("  买入价:   %.2f\n", stock.BuyPrice)
	fmt.Printf("  卖出价:   %.2f\n", stock.SellPrice)
	fmt.Println("----------------------------------------")
	fmt.Printf("  昨收价:   %.2f\n", stock.LastClose)
	fmt.Printf("  今开价:   %.2f\n", stock.OpenPrice)
	fmt.Printf("  最高价:   %.2f\n", stock.HighPrice)
	fmt.Printf("  最低价:   %.2f\n", stock.LowPrice)
	fmt.Println("----------------------------------------")
	fmt.Printf("  成交量:   %.0f\n", stock.Volume)
	fmt.Printf("  成交额:   %.2f\n", stock.Turnover)
	fmt.Printf("  时间戳:   %s\n", stock.Timestamp)
	fmt.Println("========================================")
}
