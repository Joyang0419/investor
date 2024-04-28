package history_daily_price

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"micro_stock_price/crawler"
	"repo/mongodb/schema/stock_daily_price"
	"tools/logger"
	"tools/request"
	"tools/timex"
)

type StartParamType struct {
	StockCodes        []string
	StartYearAndMonth YearAndMonth
	EndYearAndMonth   YearAndMonth
}

type crawlParamType struct {
	stockCode          string
	chooseYearAndMonth YearAndMonth
}

type CrawledDataType = []stock_daily_price.Schema

type errDataType = string

type ICommand interface {
	InsertMany(ctx context.Context, timeout time.Duration, data []stock_daily_price.Schema) (*mongo.InsertManyResult, error)
}

type Crawler struct {
	cmd ICommand
}

func NewCrawler(cmd ICommand) crawler.ICrawler[StartParamType, crawlParamType, CrawledDataType, errDataType] {
	return &Crawler{cmd: cmd}
}

func (c *Crawler) Start(startParam StartParamType, timeout time.Duration, randomSleep crawler.SleepParam) (err error) {
	if err = c.Validate(startParam); err != nil {
		return fmt.Errorf("[Crawler][Start]Validate err: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(startParam.StockCodes))
	for _, stockCode := range startParam.StockCodes {
		go func(stockCode string) {
			defer wg.Done()
			for iterator := NewYearAndMonthIterator(startParam.StartYearAndMonth, startParam.EndYearAndMonth); iterator.Next(); {
				crawled, errCrawl := c.Crawl(
					crawlParamType{
						stockCode:          stockCode,
						chooseYearAndMonth: iterator.Current(),
					},
					timeout,
					randomSleep,
				)
				if errCrawl != nil {
					if errHandleCrawlFailed := c.HandleCrawlFailed(stockCode); errHandleCrawlFailed != nil {
						logger.Error(
							"[Crawler][Start]HandleCrawlFailed stockCode:%s, err: %v",
							stockCode,
							errHandleCrawlFailed,
						)
					}
					continue
				}

				if errHandleCrawledData := c.HandleCrawledData(crawled); errHandleCrawledData != nil {
					logger.Error(
						"[Crawler][Start]HandleCrawledData stockCode:%s, err: %v",
						stockCode,
						errHandleCrawledData,
					)
				}
			}
		}(stockCode)
	}

	wg.Wait()

	return nil
}

func (c *Crawler) Validate(startParam StartParamType) error {
	if startParam.StartYearAndMonth.Year > startParam.EndYearAndMonth.Year ||
		(startParam.StartYearAndMonth.Year == startParam.EndYearAndMonth.Year && startParam.StartYearAndMonth.Month > startParam.EndYearAndMonth.Month) {
		return fmt.Errorf("[Crawler][Validate]StartYearAndMonth should be before or equal to EndYearAndMonth")
	}

	return nil
}

type twseStockData struct {
	Data [][]string `json:"data"`
}

// Crawl 爬取台灣證券交易所的股票代碼的每日數據, 來源: 台灣證券交易所個股日成交資訊
// 1次request -> 就是1個月的數據
func (c *Crawler) Crawl(crawlParam crawlParamType, timeout time.Duration, randomSleep crawler.SleepParam) (response CrawledDataType, err error) {
	logger.Info("[Crawler][Crawl]stockCode: %s, chooseYearAndMonth: %+v", crawlParam.stockCode, crawlParam.chooseYearAndMonth)

	url := fmt.Sprintf(
		"https://www.twse.com.tw/rwd/zh/afterTrading/STOCK_DAY?date=%s&stockNo=%s&response=json&_=%v",
		fmt.Sprintf("%d%02d01", crawlParam.chooseYearAndMonth.Year, crawlParam.chooseYearAndMonth.Month),
		crawlParam.stockCode,
		time.Now().UnixNano()/1000000, // 毫秒
	)

	timex.SleepRandomSeconds(randomSleep.MinSec, randomSleep.MaxSec)
	responseHttpRequest, err := request.HttpRequest[twseStockData](url, http.MethodGet, nil, timeout)
	if err != nil {
		logger.Error("[crawlTWSEDailyPrices]request.HttpRequest err: %v", err)
		return response, fmt.Errorf("[crawlTWSEDailyPrices]request.HttpRequest err: %w", err)
	}
	logger.Info("[crawlTWSEDailyPrices]received received response")

	response = make(CrawledDataType, len(responseHttpRequest.Data))
	for idx := range responseHttpRequest.Data {
		if len(responseHttpRequest.Data[idx]) != 9 {
			return nil, fmt.Errorf("len(response.Data[idx]) != 9, 被爬對象的資料格式異動了")
		}

		parsedTime, errParsedTime := parseROCDate(responseHttpRequest.Data[idx][0])
		if errParsedTime != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]timex.ParseInLocation err: %w", err)
		}

		response[idx] = stock_daily_price.Schema{}
		response[idx].DateTimestamp = parsedTime.Unix()
		response[idx].StockCode = crawlParam.stockCode

		value, errValue := strconv.ParseInt(strings.ReplaceAll(responseHttpRequest.Data[idx][1], ",", ""), 10, 64)
		if errValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseInt err: %w", errValue)
		}
		response[idx].Volume = value

		floatValue, errFloatValue := strToFloat64(responseHttpRequest.Data[idx][3])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		response[idx].OpeningPrice = floatValue

		floatValue, errFloatValue = strconv.ParseFloat(responseHttpRequest.Data[idx][4], 64)
		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		response[idx].HighestPrice = floatValue

		floatValue, errFloatValue = strToFloat64(responseHttpRequest.Data[idx][5])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		response[idx].LowestPrice = floatValue

		floatValue, errFloatValue = strToFloat64(responseHttpRequest.Data[idx][6])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		response[idx].ClosingPrice = floatValue

		// 去除空格, +, X
		re := regexp.MustCompile(`[+X]`)

		floatValue, errFloatValue = strToFloat64(
			strings.TrimSpace(
				re.ReplaceAllString(responseHttpRequest.Data[idx][7], ""),
			),
		)

		if errFloatValue != nil {
			return nil, fmt.Errorf("[crawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		response[idx].Change = floatValue
	}

	return response, nil
}

func (c *Crawler) HandleCrawledData(data CrawledDataType) error {
	_, err := c.cmd.InsertMany(context.Background(), 2*time.Second, data)
	if err != nil {
		return fmt.Errorf("[Crawler][HandleCrawledData]cmd.InsertMany err: %w", err)
	}

	return nil
}

func (c *Crawler) HandleCrawlFailed(stockCode errDataType) error {
	logger.Error("[Crawler][HandleCrawlFailed]failed stockCode: %s", stockCode)
	return nil
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
