package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	fetchAllAStocksBaseURL = "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData"
	fetchAllAStocksClient  = &http.Client{Timeout: 10 * time.Second}
	fetchAllAStocksSleep   = time.Sleep
)

// StockInfo 表示 A 股基础信息。
type StockInfo struct {
	Symbol string `json:"symbol"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Market int    `json:"market"` // 1: 沪市, 0: 深市/北交所
}

type sinaStock struct {
	Symbol string `json:"symbol"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

// FetchAllAStocks 拉取全量 A 股股票信息。
func FetchAllAStocks() ([]StockInfo, error) {
	var allStocks []StockInfo
	pageSize := 1000

	for page := 1; ; page++ {
		reqURL := fmt.Sprintf(
			"%s?page=%d&num=%d&sort=symbol&asc=1&node=hs_a&_s_r_a=setlen",
			fetchAllAStocksBaseURL,
			page,
			pageSize,
		)

		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

		resp, err := fetchAllAStocksClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("请求接口失败: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("读取响应数据失败: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("请求接口失败: status=%s body=%s", resp.Status, string(body))
		}

		pageStocks, err := parseSinaStocks(body)
		if err != nil {
			return nil, err
		}

		fmt.Println("length:", len(pageStocks))

		if len(pageStocks) == 0 {
			break
		}

		allStocks = append(allStocks, pageStocks...)

		fetchAllAStocksSleep(200 * time.Millisecond)
	}

	for _, stock := range allStocks {
		fmt.Println(stock)
	}

	return allStocks, nil
}

func parseSinaStocks(body []byte) ([]StockInfo, error) {
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" || trimmed == "null" {
		return nil, nil
	}

	var items []sinaStock
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	stocks := make([]StockInfo, 0, len(items))
	for _, item := range items {
		stocks = append(stocks, StockInfo{
			Symbol: item.Symbol,
			Code:   item.Code,
			Name:   item.Name,
			Market: marketFromSymbol(item.Symbol),
		})
	}

	return stocks, nil
}

func marketFromSymbol(symbol string) int {
	if strings.HasPrefix(symbol, "sh") {
		return 1
	}
	return 0
}
