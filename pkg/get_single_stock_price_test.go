package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStockData(t *testing.T) {
	info, err := GetStockPrice("920002")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", info)
}
