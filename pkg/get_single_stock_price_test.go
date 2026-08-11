package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStockData(t *testing.T) {
	info, err := GetStockData("600519")
	assert.Nil(t, err)
	fmt.Printf("stock info: %+v\n", info)
}
