package hisport

import (
	"context"
	"fmt"
	"time"

	"github.com/eyo-chen/expense-tracker-go/pkg/logger"
	pb "github.com/eyo-chen/expense-tracker-go/proto/hisport"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service is a struct that encapsulates the gRPC client.
type Service struct {
	client pb.HistoricalPortfolioServiceClient
}

// NewService creates a new Service instance with a gRPC client connection.
func NewService(addr string) *Service {
	conn, err := grpc.NewClient(addr, grpc.WithInsecure()) //nolint:all
	if err != nil {
		logger.Error("Failed to connect gRPC server", "error", err)
	}

	return &Service{
		client: pb.NewHistoricalPortfolioServiceClient(conn),
	}
}

func (s *Service) Create(ctx context.Context, userID int32, date time.Time) error {
	req := &pb.CreateReq{
		UserId: userID,
		Date:   timestamppb.New(date),
	}

	if _, err := s.client.Create(ctx, req); err != nil {
		logger.Error("Failed to create historical portfolio via gRPC", "error", err)
		return err
	}

	return nil
}

func (s *Service) GetPortfolioValue(ctx context.Context, userID int32, dateOption string) ([]string, []float64, error) {
	req := &pb.GetPortfolioValueReq{
		UserId:     userID,
		DateOption: dateOption,
	}

	resp, err := s.client.GetPortfolioValue(ctx, req)
	if err != nil {
		logger.Error("Failed to get portfolio value via gRPC", "error", err)
		return nil, nil, err
	}

	fmt.Println("resp.Date", resp.Date)
	fmt.Println("resp.Values", resp.Values)

	return resp.Date, resp.Values, nil
}

func (s *Service) GetGain(ctx context.Context, userID int32, dateOption string) ([]string, []float64, error) {
	req := &pb.GetGainReq{
		UserId:     userID,
		DateOption: dateOption,
	}

	resp, err := s.client.GetGain(ctx, req)
	if err != nil {
		logger.Error("Failed to get portfolio gain via gRPC", "error", err)
		return nil, nil, err
	}

	return resp.Date, resp.Values, nil
}
