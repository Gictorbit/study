package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	pb "github.com/gictorbit/study/grpc/proto/pb"
	server "github.com/gictorbit/study/grpc/server"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	serverHost string = "0.0.0.0"
	serverPort int    = 3456
)

func main() {
	logger := log.Default()
	helloServer := server.NewHelloServer(logger)
	goodByeServer := server.NewGoodByeServer(logger)

	unaryOpts := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
	}
	streamOpts := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamOpts...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryOpts...)),
	)
	socket := fmt.Sprintf("%s:%d", serverHost, serverPort)
	lis, err := net.Listen("tcp", socket)
	if err != nil {
		logger.Fatalf("error listerner: %v\n", err)
	}

	reflection.Register(s)
	pb.RegisterHelloServiceServer(s, helloServer)
	pb.RegisterGoodByeServiceServer(s, goodByeServer)

	go func() {
		err := s.Serve(lis)
		if err != nil {
			logger.Fatalf("error listerner: %v\n", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	logger.Printf("server started on %s\n", socket)
	// Block until a signal is received
	<-ch
	logger.Println("exit grpc server")
	s.Stop()
	lis.Close()
}
