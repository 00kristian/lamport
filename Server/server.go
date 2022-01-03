package main

import (
	"context"
	"log"
	"net"
	"sync"

	"GRPC_demo/chat"

	"google.golang.org/grpc"
)

type SafeLamport struct {
	mu sync.Mutex
	Lamport int64
	}
	
func (c *SafeLamport) Tick() {
	c.mu.Lock()
	c.Lamport++
	defer c.mu.Unlock()
}

func(c *SafeLamport) CompareClock(otherClockTime int64) {
	c.mu.Lock()
	newVal :=  max(c.Lamport, otherClockTime) 
	c.Lamport = newVal
	defer c.mu.Unlock()
}

type Server struct {
	chat.UnimplementedChatServiceServer
	SafeLamport SafeLamport
}

func main() {
	lis, lisErr := net.Listen("tcp", ":9000")
	if lisErr != nil {
		log.Fatalf("Failed to listen on port 9000: %v", lisErr)
	}

	s := Server{SafeLamport: SafeLamport{Lamport: 0}}

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)

	serveErr := grpcServer.Serve(lis)
	if serveErr != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", serveErr)
	}


}


func (s *Server) SayHello(ctx context.Context, message *chat.Message) (*chat.Message, error) { 

	log.Printf("Recieved Message body from client: %s lamport: %d" , message.Body, s.SafeLamport.Lamport)

//message.Lamport
	s.SafeLamport.Tick()
	return &chat.Message{Body: "Hello from the server! wit timestamp", Lamport:s.SafeLamport.Lamport  }, nil
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}