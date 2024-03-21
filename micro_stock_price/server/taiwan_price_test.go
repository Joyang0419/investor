package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCrawlTWSEDailyPrices(t *testing.T) {
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		fmt.Println("加载时区失败:", err)
		return
	}
	prices, err := CrawlTWSEDailyPrices(
		"2330",
		time.Date(2023, time.January, 1, 0, 0, 0, 0, location),
	)

	assert.NotEmpty(t, prices)
	assert.NoError(t, err)
}
