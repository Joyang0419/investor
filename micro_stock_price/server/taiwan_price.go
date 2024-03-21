package server

import (
	"context"

	"protos/micro_stock_price"
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
