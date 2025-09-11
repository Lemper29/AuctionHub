package server

import (
	"context"
	pb "github/auction/auction-service/gen/proto"
	"github/auction/auction-service/internal/service"
	"github/auction/auction-service/internal/storage"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuctionServiceServer
	addr    string
	service *service.LotService
}

func NewGrpcServer(addr string, storage storage.Storage) *server {
	return &server{
		addr:    addr,
		service: service.NewLotService(storage),
	}
}

func (s *server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuctionServiceServer(grpcServer, s)

	log.Printf("Server starting on %s", s.addr)
	return grpcServer.Serve(lis)
}

// Реализации gRPC методов
func (s *server) CreateLot(ctx context.Context, req *pb.CreateLotRequest) (*pb.CreateLotResponse, error) {
	return s.service.CreateLot(ctx, req)
}

func (s *server) GetLot(ctx context.Context, req *pb.GetLotRequest) (*pb.GetLotResponse, error) {
	return s.service.GetLot(ctx, req)
}

func (s *server) PlaceBid(ctx context.Context, req *pb.PlaceBidRequest) (*pb.PlaceBidResponse, error) {
	return s.service.PlaceBid(ctx, req)
}

func (s *server) SubscribeToLot(req *pb.SubscribeToLotRequest, stream pb.AuctionService_SubscribeToLotServer) error {
	return s.service.SubscribeToLot(req, stream)
}
