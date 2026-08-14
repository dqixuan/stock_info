package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStockMarginByDate(t *testing.T) {
	info, err := GetStockMargin("920000")
	assert.Nil(t, err)
	fmt.Printf("stock margin: %+v", info)
}
