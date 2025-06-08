package stock

import (
	"context"

	"github.com/eyo-chen/expense-tracker-go/internal/domain"
	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	pb "github.com/eyo-chen/expense-tracker-go/proto/stock"
	"google.golang.org/grpc"
)

// Service is a struct that encapsulates the gRPC client.
type Service struct {
	client pb.StockServiceClient
}

// NewService creates a new Service instance with a gRPC client connection.
func NewService(addr string) *Service {
	conn, err := grpc.NewClient(addr)
	if err != nil {
		logger.Error("Failed to connect gRPC server", "error", err)
	}

	return &Service{
		client: pb.NewStockServiceClient(conn),
	}
}

func (s *Service) Create(ctx context.Context, stock domain.CreateStock) (string, error) {
	req := &pb.CreateReq{
		UserId:    stock.UserID,
		Symbol:    stock.Symbol,
		Price:     stock.Price,
		Quantity:  stock.Quantity,
		Action:    pb.Action_Type(pb.Action_Type_value[string(stock.ActionType)]),
		StockType: pb.StockType_Type(pb.StockType_Type_value[string(stock.StockType)]),
	}

	resp, err := s.client.Create(ctx, req)
	if err != nil {
		logger.Error("Failed to create stock via gRPC", "error", err)
		return "", err
	}

	return resp.GetId(), nil
}

func (s *Service) GetPortfolioInfo(ctx context.Context, userID int32) (domain.Portfolio, error) {
	resp, err := s.client.GetPortfolioInfo(ctx, &pb.GetPortfolioInfoReq{
		UserId: userID,
	})
	if err != nil {
		logger.Error("Failed to get portfolio info via gRPC", "error", err)
		return domain.Portfolio{}, err
	}

	return domain.Portfolio{
		UserID:              resp.GetUserId(),
		TotalPortfolioValue: resp.GetTotalPortfolioValue(),
		TotalGain:           resp.GetTotalGain(),
		ROI:                 resp.GetRoi(),
	}, nil
}

func (s *Service) GetStockInfo(ctx context.Context, userID int32) (domain.AllStockInfo, error) {
	resp, err := s.client.GetStockInfo(ctx, &pb.GetStockInfoReq{
		UserId: userID,
	})
	if err != nil {
		logger.Error("Failed to get stock info via gRPC", "error", err)
		return domain.AllStockInfo{}, err
	}

	return domain.AllStockInfo{
		Stocks: cvtToStockInfoList(resp.GetStocks()),
		ETF:    cvtToStockInfoList(resp.GetEtf()),
		Cash:   cvtToStockInfoList(resp.GetCash()),
	}, nil
}

func cvtToStockInfoList(stocks []*pb.StockInfo) []domain.StockInfo {
	result := make([]domain.StockInfo, len(stocks))
	for i, s := range stocks {
		result[i] = domain.StockInfo{
			Symbol:     s.GetSymbol(),
			Quantity:   s.GetQuantity(),
			Price:      s.GetPrice(),
			AvgCost:    s.GetAvgCost(),
			Percentage: s.GetPercentage(),
		}
	}

	return result
}
