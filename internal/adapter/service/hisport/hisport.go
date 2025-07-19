package hisport

import (
	"context"
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
		logger.Error("Failed to create stock via gRPC", "error", err)
		return err
	}

	return nil
}
