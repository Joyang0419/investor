package taiwan_stock

import (
	"context"

	"protos/micro_stock_price"
)

type TaiwanStockServer struct {
	micro_stock_price.UnimplementedTaiwanPriceServer
}

func NewTaiwanPriceServer() micro_stock_price.TaiwanPriceServer {
	return &TaiwanStockServer{}
}

func (server *TaiwanStockServer) GetDailyPrices(_ context.Context, _ *micro_stock_price.DailyPricesRequest) (*micro_stock_price.DailyPricesResponse, error) {
	//TODO implement me
	panic("implement me")
}
