package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"protos/micro_stock_price"
	"tools/request"
)

type TaiwanPriceServer struct {
	micro_stock_price.UnimplementedTaiwanPriceServer
}

func NewTaiwanPriceServer() micro_stock_price.TaiwanPriceServer {
	return new(TaiwanPriceServer)
}

// todo README.md
// 目前覺得一層就好(trace code 方便), 因為gateway server 已經處理好參數 丟下來了, 之後 repo 裝在封裝另一個個地方，這邊 就直接把邏輯寫在這，要小封裝，就直接用方法  在下面的檔案用func吧

func (server *TaiwanPriceServer) GetDailyPrices(ctx context.Context, request *micro_stock_price.DailyPricesRequest) (*micro_stock_price.DailyPricesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (server *TaiwanPriceServer) CrawlDailyPrices(ctx context.Context, request *micro_stock_price.DailyPricesRequest) (*micro_stock_price.DailyPricesResponse, error) {
	//TODO implement me
	panic("implement me")
}

// TWSEStockDailyPrice 台灣證券交易所個股日成交資訊
type TWSEStockDailyPrice struct {
	Date         string  // 日期
	Volume       int64   // 成交股數
	Transaction  int64   // 成交金額
	OpeningPrice float64 // 開盤價
	HighestPrice float64 // 最高價
	LowestPrice  float64 // 最低價
	ClosingPrice float64 // 收盤價
	Change       float64 // 漲跌價差
	Transactions int64   // 成交筆數
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

	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return nil, fmt.Errorf("[CrawlTWSEDailyPrices]time.LoadLocation err: %w", err)
	}

	prices = make([]*micro_stock_price.DailyPrice, len(prices))
	for idx := range response.Data {
		if len(response.Data[idx]) != 9 {
			return nil, fmt.Errorf("len(response.Data[idx]) != 9, 被爬對象的資料格式異動了")
		}

		// FIXME 時間是中華民國時間
		parsedTime, errParsedTime := time.ParseInLocation("2024/03/01", response.Data[idx][0], location)
		if errParsedTime != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]time.ParseInLocation err: %w", err)
		}

		// FIXME 15,311,364   因為有, Parse不出來
		prices[idx].DateTimestamp = parsedTime.Unix()
		prices[idx].StockCode = stockCode

		a := strings.ReplaceAll(response.Data[idx][1], ",", "")
		_ = a
		value, errValue := strconv.ParseInt(strings.ReplaceAll(response.Data[idx][1], ",", ""), 10, 64)
		if errValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseInt err: %w", errValue)
		}
		prices[idx].Volume = value

		floatValue, errFloatValue := strconv.ParseFloat(response.Data[idx][3], 64)
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].OpeningPrice = floatValue

		floatValue, errFloatValue = strconv.ParseFloat(response.Data[idx][4], 64)
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].HighestPrice = floatValue

		floatValue, errFloatValue = strconv.ParseFloat(response.Data[idx][5], 64)
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].LowestPrice = floatValue

		floatValue, errFloatValue = strconv.ParseFloat(response.Data[idx][6], 64)
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].ClosingPrice = floatValue

		floatValue, errFloatValue = strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(response.Data[idx][7], "+", ""), "-", ""), 64)
		if errFloatValue != nil {
			return nil, fmt.Errorf("[CrawlTWSEDailyPrices]strconv.ParseFloat err: %w", errFloatValue)
		}
		prices[idx].Change = floatValue
	}

	return prices, nil
}
