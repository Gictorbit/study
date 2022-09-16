package server

import (
	"context"
	"fmt"
	"log"

	pb "github.com/gictorbit/study/grpc/proto/pb"
)

type MainServer struct {
	pb.UnimplementedMainServiceServer
	*HelloServer
	*GoodByeServer
}

type HelloServer struct {
	names  map[string]int
	logger *log.Logger
	pb.UnimplementedMainServiceServer
}

type GoodByeServer struct {
	names  map[string]int
	logger *log.Logger
	pb.UnimplementedMainServiceServer
}

func NewHelloServer(logger *log.Logger) *HelloServer {
	return &HelloServer{
		logger: logger,
		names:  make(map[string]int),
	}
}

func NewGoodByeServer(logger *log.Logger) *GoodByeServer {
	return &GoodByeServer{
		logger: logger,
		names:  make(map[string]int),
	}
}

func NewMainServer(logger *log.Logger) *MainServer {
	return &MainServer{
		HelloServer:   NewHelloServer(logger),
		GoodByeServer: NewGoodByeServer(logger),
	}
}

func (h *HelloServer) SayHello(ctx context.Context, req *pb.SayHelloRequest) (*pb.SayHelloResponse, error) {

	h.names[req.Name]++
	msg := fmt.Sprintf("Hello %s", req.Name)

	h.logger.Printf("SayHello: %s", req.Name)
	h.logger.Println("msg: ", msg)

	return &pb.SayHelloResponse{
		Msg: msg,
	}, nil
}

func (g *GoodByeServer) SayGoodBye(ctx context.Context, req *pb.SayGoodByeRequest) (*pb.SayGoodByeResponse, error) {

	g.names[req.Name]++
	msg := fmt.Sprintf("GoodBye %s", req.Name)

	g.logger.Printf("SayGoodBye: %s", req.Name)
	g.logger.Println("msg: ", msg)

	return &pb.SayGoodByeResponse{
		Msg: msg,
	}, nil
}

func (m *MainServer) SayHello(ctx context.Context, req *pb.SayHelloRequest) (*pb.SayHelloResponse, error) {
	return m.HelloServer.SayHello(ctx, req)
}

func (m *MainServer) SayGoodBye(ctx context.Context, req *pb.SayGoodByeRequest) (*pb.SayGoodByeResponse, error) {
	return m.GoodByeServer.SayGoodBye(ctx, req)
}
