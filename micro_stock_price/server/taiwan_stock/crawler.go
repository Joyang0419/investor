package taiwan_stock

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"repo/mysql/schema"
	"tools/logger"
	"tools/request"
	"tools/timex"
)

type ITaiwanStockCrawler interface {
	DailyPrices(
		ctx context.Context,
		stockCodes []string,
		startYearAndMonth YearAndMonth,
		endYearAndMonth YearAndMonth,
	) (
		ch chan []schema.DailyPrice,
		errChStockCode chan string,
		errStockCodes []string,
	)
}

type TaiwanStockCrawler struct{}

func NewTaiwanStockCrawler() *TaiwanStockCrawler {
	return &TaiwanStockCrawler{}
}

type YearAndMonth struct {
	Year  uint32
	Month uint32
}

func (crawler *TaiwanStockCrawler) DailyPrices(
	_ context.Context,
	stockCodes []string,
	startYearAndMonth YearAndMonth,
	endYearAndMonth YearAndMonth,
) (
	ch chan []schema.DailyPrice,
	errChStockCode chan string,
	err error,
) {
	// 檢查 startYearAndMonth 和 endYearAndMonth 是否合法
	if err = ValidateYearAndMonth(startYearAndMonth, endYearAndMonth); err != nil {
		return nil, nil, fmt.Errorf("[DailyPrices]ValidateYearAndMonth err: %w", err)
	}

	ch = make(chan []schema.DailyPrice, len(stockCodes))
	errChStockCode = make(chan string, len(stockCodes))

	var wg sync.WaitGroup
	wg.Add(len(stockCodes))
	for _, stockCode := range stockCodes {
		go func(stockCode string) {
			defer wg.Done()
			for iterator := NewYearAndMonthIterator(startYearAndMonth, endYearAndMonth); iterator.Next(); {
				prices, errPrices := crawlTWSEDailyPrices(stockCode, iterator.Current())
				if errPrices != nil {
					errChStockCode <- stockCode
					continue
				}
				ch <- prices
			}
		}(stockCode)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(errChStockCode)
	}()

	return ch, errChStockCode, nil
}

type twseStockData struct {
	Data [][]string `json:"data"`
}

func ValidateYearAndMonth(startYearAndMonth, endYearAndMonth YearAndMonth) error { // Check if startYearAndMonth is after endYearAndMonth
	if startYearAndMonth.Year > endYearAndMonth.Year ||
		(startYearAndMonth.Year == endYearAndMonth.Year && startYearAndMonth.Month > endYearAndMonth.Month) {
		return fmt.Errorf("startYearAndMonth should be before or equal to endYearAndMonth")
	}

	return nil
}

// crawlTWSEDailyPrices 爬取台灣證券交易所的股票代碼的每日數據, 來源: 台灣證券交易所個股日成交資訊
// 1次request -> 就是1個月的數據
func crawlTWSEDailyPrices(stockCode string, chooseYearAndMonth YearAndMonth) (prices []schema.DailyPrice, err error) {
	logger.Info("[crawlTWSEDailyPrices]stockCode: %s, chooseYearAndMonth: %+v", stockCode, chooseYearAndMonth)
	url := fmt.Sprintf(
		"https://www.twse.com.tw/rwd/zh/afterTrading/STOCK_DAY?date=%s&stockNo=%s&response=json&_=%v",
		fmt.Sprintf("%d%02d01", chooseYearAndMonth.Year, chooseYearAndMonth.Month),
		stockCode,
		time.Now().UnixNano()/1000000, // 毫秒
	)
	// 隨機等待1~5秒
	timex.SleepRandomSeconds(1, 5)
	response, err := request.HttpRequest[twseStockData](url, http.MethodGet, nil, 10*time.Second)
	if err != nil {
		logger.Error("[crawlTWSEDailyPrices]request.HttpRequest err: %v", err)
		return nil, fmt.Errorf("[crawlTWSEDailyPrices]request.HttpRequest err: %w", err)
	}
	logger.Info("[crawlTWSEDailyPrices]received received response")

	prices = make([]schema.DailyPrice, len(response.Data))
	for idx := range response.Data {
		if len(response.Data[idx]) != 9 {
			return nil, fmt.Errorf("len(response.Data[idx]) != 9, 被爬對象的資料格式異動了")
		}

		parsedTime, errParsedTime := parseROCDate(response.Data[idx][0])
		if errParsedTime != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]timex.ParseInLocation err: %w", err)
		}

		prices[idx] = schema.DailyPrice{}
		prices[idx].DateTimestamp = parsedTime.Unix()
		prices[idx].StockCode = stockCode

		value, errValue := strconv.ParseInt(strings.ReplaceAll(response.Data[idx][1], ",", ""), 10, 64)
		if errValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseInt err: %w", errValue)
		}
		prices[idx].Volume = value

		floatValue, errFloatValue := strToFloat64(response.Data[idx][3])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].OpeningPrice = floatValue

		floatValue, errFloatValue = strconv.ParseFloat(response.Data[idx][4], 64)
		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].HighestPrice = floatValue

		floatValue, errFloatValue = strToFloat64(response.Data[idx][5])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].LowestPrice = floatValue

		floatValue, errFloatValue = strToFloat64(response.Data[idx][6])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].ClosingPrice = floatValue

		// 去除空格, +, X
		re := regexp.MustCompile(`[+X]`)

		floatValue, errFloatValue = strToFloat64(
			strings.TrimSpace(
				re.ReplaceAllString(response.Data[idx][7], ""),
			),
		)

		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].Change = floatValue
	}

	return prices, nil
}

// parseROCDate 將中華民國紀年的日期字符串轉換為time.Time ex.112/01/03
func parseROCDate(rocDateStr string) (time.Time, error) {
	// 民國年與西元年的轉換基數
	const rocToAD = 1911

	// 將日期字符串分割為年、月、日
	parts := strings.Split(rocDateStr, "/")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid ROC date format")
	}

	// 轉換年份
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year: %s", parts[0])
	}
	year += rocToAD // 將民國年轉為西元年

	// 重組為西元年的日期字符串
	adDateStr := fmt.Sprintf("%d/%s/%s", year, parts[1], parts[2])

	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return time.Time{}, fmt.Errorf("[parseROCDate]timex.LoadLocation err: %w", err)
	}

	// 解析日期
	t, err := time.ParseInLocation("2006/01/02", adDateStr, location)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func strToFloat64(s string) (float64, error) {
	if s == "0.00" {
		return 0, nil
	}

	return strconv.ParseFloat(s, 64)
}
