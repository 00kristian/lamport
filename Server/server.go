package main

import (
	"context"
	"log"
	"net"

	"GRPC_demo/chat"

	"google.golang.org/grpc"
)

type Server struct {
	chat.UnimplementedChatServiceServer
}

func main() {
	lis, lisErr := net.Listen("tcp", ":9000")
	if lisErr != nil {
		log.Fatalf("Failed to listen on port 9000: %v", lisErr)
	}

	s := Server{chat.UnimplementedChatServiceServer{}}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)

	serveErr := grpcServer.Serve(lis)
	if serveErr != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", serveErr)
	}

}


func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Recieved Message body from client: %s", message.Body)
	return &Message{Body: "Hello from the server!"}, nil
}