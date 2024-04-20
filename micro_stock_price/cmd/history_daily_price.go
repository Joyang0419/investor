package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"micro_stock_price/crawler"
	"micro_stock_price/crawler/history_daily_price"
	"repo/mongodb/stock_daily_price"
	"tools/infra_conn"
	"tools/logger"
)

var dailyPriceCmd = &cobra.Command{
	Use:   "dailyPrice",
	Short: "",
	Long:  "",
	Run:   runDailyPriceCmd,
}

func init() {
	rootCmd.AddCommand(dailyPriceCmd)
}

func runDailyPriceCmd(_ *cobra.Command, _ []string) {
	mongoDbConn, err := infra_conn.SetupMongoDB(
		infra_conn.MongoDBCfg{
			Host:            "127.0.0.1",
			Port:            27017,
			Username:        "root",
			Password:        "root",
			Database:        "admin",
			ConnectTimeout:  20 * time.Second,
			MaxPoolSize:     20,
			MaxConnIdleTime: 15 * time.Minute,
		},
	)
	if err != nil {
		logger.Fatal("[runDailyPriceCmd]infra_conn.SetupMongoDB err: %v", err)
	}

	// 設定關鍵參數
	// 1. 股票代碼(可複數)
	stockCodes := []string{"2330", "2317"}
	// 2. 開始時間
	startYearAndMonth := history_daily_price.YearAndMonth{
		Year:  2021,
		Month: 1,
	}
	// 3. 結束時間
	endYearAndMonth := history_daily_price.YearAndMonth{
		Year:  2021,
		Month: 3,
	}

	logger.Info(
		"stockCodes: %+v, startYearAndMonth: %+v, endYearAndMonth: %+v",
		stockCodes,
		startYearAndMonth,
		endYearAndMonth,
	)

	stockDailyPriceCmd := stock_daily_price.NewCommand(mongoDbConn)
	c := history_daily_price.NewCrawler(stockDailyPriceCmd)

	if err = c.Start(
		history_daily_price.StartParamType{
			StockCodes:        stockCodes,
			StartYearAndMonth: startYearAndMonth,
			EndYearAndMonth:   endYearAndMonth,
		},
		5*time.Second,
		crawler.SleepParam{
			MinSec: 1,
			MaxSec: 3,
		},
	); err != nil {
		logger.Fatal("[runDailyPriceCmd]history_daily_price.Crawler.Start err: %v", err)
	}

	logger.Info("[runDailyPriceCmd]CrawlDailyPrices done")
}
