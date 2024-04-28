package taiwan_stock

import (
	"context"

	"protos/micro_stock_price"
)

type Server struct {
	micro_stock_price.UnimplementedTaiwanPriceServer
}

func NewServer() micro_stock_price.TaiwanPriceServer {
	return &Server{}
}

func (server *Server) GetDailyPrices(_ context.Context, _ *micro_stock_price.DailyPricesRequest) (*micro_stock_price.DailyPricesResponse, error) {
	//TODO implement me
	panic("implement me")
}
