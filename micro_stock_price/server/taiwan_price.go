package server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"protos/micro_stock_price"
	"tools/request"
	"tools/timex"
)

type TaiwanPriceServer struct {
	micro_stock_price.UnimplementedTaiwanPriceServer
}

func NewTaiwanPriceServer() micro_stock_price.TaiwanPriceServer {
	return new(TaiwanPriceServer)
}

func (server *TaiwanPriceServer) GetDailyPrices(ctx context.Context, request *micro_stock_price.DailyPricesRequest) (*micro_stock_price.DailyPricesResponse, error) {
	//TODO implement me
	panic("implement me")
}

// CrawlDailyPrices 將startTime 和 endTime  切成 2022/01/01, 2022/02/01 ... 再往下丟 CrawlTWSEDailyPrices 去TWSE 拿資料(他打資料，用月份為單位去拿)
// startTimeStamp example: 1633039200
func (server *TaiwanPriceServer) CrawlDailyPrices(_ context.Context, request *micro_stock_price.DailyPricesRequest) (*micro_stock_price.DailyPricesResponse, error) {
	startTime := timex.Int64ToTime(request.StartTimestamp)
	endTime := timex.Int64ToTime(request.EndTimestamp)

	var allPrices []*micro_stock_price.DailyPrice
	// 遍歷每個月份
	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.AddDate(0, 1, 0) {
		prices, err := CrawlTWSEDailyPrices(request.StockCode, t)
		if err != nil {
			return nil, fmt.Errorf("[TaiwanPriceServer][CrawlDailyPrices]CrawlTWSEDailyPrices err: %w", err)
		}

		allPrices = append(allPrices, prices...)
	}

	return &micro_stock_price.DailyPricesResponse{DailyPrices: allPrices}, nil
}

type TWSEStockData struct {
	Data [][]string `json:"data"`
}

// CrawlTWSEDailyPrices 爬取台灣證券交易所的股票代碼的每日數據, 來源: 台灣證券交易所個股日成交資訊
// 1次request -> 就是1個月的數據
func CrawlTWSEDailyPrices(stockCode string, chooseDateTime time.Time) (prices []*micro_stock_price.DailyPrice, err error) {
	url := fmt.Sprintf(
		"https://www.twse.com.tw/rwd/zh/afterTrading/STOCK_DAY?date=%s&stockNo=%s&response=json&_=%v",
		chooseDateTime.Format("20060102"),
		stockCode,
		time.Now().UnixNano()/1000000, // 毫秒
	)
	response, err := request.HttpRequest[TWSEStockData](url, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("[CrawlTWSEDailyPrices]request.HttpRequest err: %w", err)
	}

	prices = make([]*micro_stock_price.DailyPrice, len(response.Data))
	for idx := range response.Data {
		if len(response.Data[idx]) != 9 {
			return nil, fmt.Errorf("len(response.Data[idx]) != 9, 被爬對象的資料格式異動了")
		}

		parsedTime, errParsedTime := ParseROCDate(response.Data[idx][0])
		if errParsedTime != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]timex.ParseInLocation err: %w", err)
		}

		prices[idx] = new(micro_stock_price.DailyPrice)
		prices[idx].DateTimestamp = parsedTime.Unix()
		prices[idx].StockCode = stockCode

		value, errValue := strconv.ParseInt(strings.ReplaceAll(response.Data[idx][1], ",", ""), 10, 64)
		if errValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseInt err: %w", errValue)
		}
		prices[idx].Volume = value

		floatValue, errFloatValue := StrToFloat64(response.Data[idx][3])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].OpeningPrice = floatValue

		floatValue, errFloatValue = strconv.ParseFloat(response.Data[idx][4], 64)
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].HighestPrice = floatValue

		floatValue, errFloatValue = StrToFloat64(response.Data[idx][5])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].LowestPrice = floatValue

		floatValue, errFloatValue = StrToFloat64(response.Data[idx][6])
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].ClosingPrice = floatValue

		// 去除空格, +, X
		re := regexp.MustCompile(`[\+X]`)

		floatValue, errFloatValue = StrToFloat64(
			strings.TrimSpace(
				re.ReplaceAllString(response.Data[idx][7], ""),
			),
		)

		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].Change = floatValue
	}

	return prices, nil
}

// ParseROCDate 將中華民國紀年的日期字符串轉換為time.Time ex.112/01/03
func ParseROCDate(rocDateStr string) (time.Time, error) {
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
		return time.Time{}, fmt.Errorf("[ParseROCDate]timex.LoadLocation err: %w", err)
	}

	// 解析日期
	t, err := time.ParseInLocation("2006/01/02", adDateStr, location)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func StrToFloat64(s string) (float64, error) {
	if s == "0.00" {
		return 0, nil
	}

	return strconv.ParseFloat(s, 64)
}
