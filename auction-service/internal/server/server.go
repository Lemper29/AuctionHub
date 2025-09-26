package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/Lemper29/auction-service/internal/service"
	"github.com/Lemper29/auction-service/internal/storage"
	pb "github.com/Lemper29/auction/gen/auction"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuctionServiceServer
	addr    string
	service *service.LotService
	logger  *slog.Logger
}

func NewGrpcServer(addr string, storage storage.Storage, appLogger *slog.Logger) *server {
	serverLogger := appLogger.With("component", "grpc-server")

	return &server{
		addr:    addr,
		service: service.NewLotService(storage, serverLogger),
		logger:  serverLogger,
	}
}

func (s *server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.ErrorContext(context.Background(), "Failed to listen", "address", s.addr, "error", err)
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuctionServiceServer(grpcServer, s)

	s.logger.InfoContext(context.Background(), "Server starting", "address", s.addr)

	if err := grpcServer.Serve(lis); err != nil {
		s.logger.ErrorContext(context.Background(), "Server failed", "error", err)
		return err
	}

	return nil
}

// Реализации gRPC методов
func (s *server) CreateLot(ctx context.Context, req *pb.CreateLotRequest) (*pb.CreateLotResponse, error) {
	s.logger.DebugContext(ctx, "CreateLot called", "name", req.Name)
	return s.service.CreateLot(ctx, req)
}

func (s *server) GetLot(ctx context.Context, req *pb.GetLotRequest) (*pb.GetLotResponse, error) {
	s.logger.DebugContext(ctx, "GetLot called", "lot_id", req.LotId)
	return s.service.GetLot(ctx, req)
}

func (s *server) PlaceBid(ctx context.Context, req *pb.PlaceBidRequest) (*pb.PlaceBidResponse, error) {
	s.logger.DebugContext(ctx, "PlaceBid called",
		"lot_id", req.LotId,
		"user_id", req.UserId,
		"amount", req.Amount,
	)
	return s.service.PlaceBid(ctx, req)
}

func (s *server) SubscribeToLot(req *pb.SubscribeToLotRequest, stream pb.AuctionService_SubscribeToLotServer) error {
	s.logger.InfoContext(stream.Context(), "SubscribeToLot called", "lot_id", req.LotId)
	return s.service.SubscribeToLot(req, stream)
}
