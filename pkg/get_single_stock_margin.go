package pkg

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

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

type marginResponse struct {
	Error string          `json:"error"`
	Data  json.RawMessage `json:"data"`
}

type marginPayload struct {
	Date               string  `json:"信用交易日期"`
	MarginBalance      float64 `json:"融资余额"`
	MarginBuyAmount    float64 `json:"融资买入额"`
	MarginRepayAmount  float64 `json:"融资偿还额"`
	ShortSellVolume    float64 `json:"融券余量"`
	ShortSellAmount    float64 `json:"融券卖出量"`
	ShortRepayAmount   float64 `json:"融券偿还量"`
	TotalMarginBalance float64 `json:"融资融券余额"`
	ShortBalance       float64 `json:"融券余额"`
}

func GetStockMargin(symbol string) (*MarginData, error) {
	return GetStockMarginByDate(symbol, "")
}

func GetStockMarginByDate(symbol, tradeDate string) (*MarginData, error) {
	scriptPath, err := resolveScriptPath("get_stock_margin_info.py")
	if err != nil {
		return nil, err
	}

	// 调用 Python 脚本，传入股票代码
	args := []string{scriptPath, symbol}
	if tradeDate != "" {
		args = append(args, tradeDate, tradeDate)
	}
	cmd := exec.Command("python3", args...)

	// 执行命令并获取标准输出
	output, err := cmd.Output()
	if err != nil {
		// 捕获 Python 脚本的错误输出（stderr）
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("脚本执行失败: %s\nStderr: %s", err.Error(), string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("脚本执行失败: %w", err)
	}

	var resp marginResponse
	if err = json.Unmarshal(output, &resp); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("脚本返回错误: %s", resp.Error)
	}
	if len(resp.Data) == 0 || string(resp.Data) == "null" {
		return nil, fmt.Errorf("脚本未返回融资融券数据")
	}

	payload, err := parseMarginPayload(resp.Data)
	if err != nil {
		return nil, err
	}

	return &MarginData{
		Date:               payload.Date,
		MarginBalance:      payload.MarginBalance,
		MarginBuyAmount:    payload.MarginBuyAmount,
		MarginRepayAmount:  payload.MarginRepayAmount,
		ShortSellVolume:    payload.ShortSellVolume,
		ShortSellAmount:    payload.ShortSellAmount,
		ShortRepayAmount:   payload.ShortRepayAmount,
		TotalMarginBalance: payload.TotalMarginBalance,
		ShortBalance:       payload.ShortBalance,
	}, nil
}

func parseMarginPayload(raw json.RawMessage) (*marginPayload, error) {
	var single marginPayload
	if err := json.Unmarshal(raw, &single); err == nil && single.Date != "" {
		return &single, nil
	}

	var items []marginPayload
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("融资融券数据解析失败: %w", err)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("融资融券数据为空")
	}

	latest := items[0]
	for _, item := range items[1:] {
		if item.Date > latest.Date {
			latest = item
		}
	}
	return &latest, nil
}
