package pkg

import (
	"encoding/json"
	"fmt"
	"os/exec"
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

// MarginData 结构体，用于解析融资融券数据
type MarginData struct {
	Date               string  `json:"日期"`
	MarginBalance      float64 `json:"融资余额"`
	MarginBuyAmount    float64 `json:"融资买入额"`
	MarginRepayAmount  float64 `json:"融资偿还额"`
	ShortSellVolume    float64 `json:"融券余量"`
	ShortSellAmount    float64 `json:"融券卖出量"`
	ShortRepayAmount   float64 `json:"融券偿还量"`
	TotalMarginBalance float64 `json:"融资融券余额"`
	ShortBalance       float64 `json:"融券余额"`
}

func GetStockData(symbol string) (*StockData, error) {
	// 调用 Python 脚本，传入股票代码
	cmd := exec.Command("python3", "../scripts/get_stock_price.py", symbol)

	// 执行命令并获取标准输出
	output, err := cmd.Output()
	if err != nil {
		// 捕获 Python 脚本的错误输出（stderr）
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("脚本执行失败: %s\nStderr: %s", err.Error(), string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("脚本执行失败: %w", err)
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
