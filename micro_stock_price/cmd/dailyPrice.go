package cmd

import (
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"micro_stock_price/server/taiwan_stock"
	"repo/mysql"
	"tools/infra_conn"
	"tools/logger"
	"tools/slicex"
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

func runDailyPriceCmd(cmd *cobra.Command, _ []string) {
	mysqlConn, err := infra_conn.SetupMySQL(
		infra_conn.MySQLCfg{
			Host:            "localhost",
			Port:            3306,
			Username:        "root",
			Password:        "root",
			Database:        "investor",
			MaxIdleConns:    10,
			MaxOpenConns:    10,
			ConnMaxLifeTime: 60 * time.Second,
		}, logger.GormInfoLogger(),
	)
	if err != nil {
		log.Fatalf("infra_conn.SetupMySQL err: %v", err)
	}

	repo := mysql.NewTaiwanStockRepo(mysqlConn)
	crawler := taiwan_stock.NewTaiwanStockCrawler()

	// 設定關鍵參數
	// 1. 股票代碼(可複數)
	stockCodes := []string{"2330", "2317"}
	// 2. 開始時間
	startYearAndMonth := taiwan_stock.YearAndMonth{
		Year:  2021,
		Month: 1,
	}
	// 3. 結束時間
	endYearAndMonth := taiwan_stock.YearAndMonth{
		Year:  2021,
		Month: 3,
	}

	logger.Info(
		"stockCodes: %+v, startYearAndMonth: %+v, endYearAndMonth: %+v",
		stockCodes,
		startYearAndMonth,
		endYearAndMonth,
	)
	pricesCh, errChStock, err := crawler.DailyPrices(
		cmd.Context(),
		stockCodes,
		startYearAndMonth,
		endYearAndMonth,
	)
	if err != nil {
		logger.Fatal("crawler.DailyPrices err: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var errStockCodes []string
	go func() {
		defer wg.Done()
		for prices := range pricesCh {
			if err = repo.BulkInsertDailyPrices(prices); err != nil {
				if slicex.IsNotEmpty(prices) {
					errStockCodes = append(errStockCodes, prices[0].StockCode)
				}
			}
		}
	}()
	go func() {
		defer wg.Done()
		for errStockCode := range errChStock {
			errStockCodes = append(errStockCodes, errStockCode)
		}
	}()
	wg.Wait()
	if slicex.IsNotEmpty(errStockCodes) {
		// 準備手動執行從跑的 stockCode && 記得刪除資料庫既有的資料，直接重跑
		logger.Fatal("[TaiwanStockServer][DailyPrices]failed stockCodes: %v", errStockCodes)
	}
}
