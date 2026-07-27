package pkg

import (
	"fmt"
	"testing"
)

func TestFetchAllAStocks(t *testing.T) {

	stocks, err := FetchAllAStocks()
	if err != nil {
		t.Fatalf("FetchAllAStocks() error = %v", err)
	}

	for _, stock := range stocks {
		fmt.Println(stock)
	}
}
