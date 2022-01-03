package main

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"GRPC_demo/chat"
)



var safeL = SafeLamport{}

type SafeLamport struct {
	mu sync.Mutex
	Lamport int64
	}


func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)
	for{
	safeL.tickmychild(context.Background() ,c)
	time.Sleep(3*time.Second)
	}


}


func (s *SafeLamport) tickmychild (ctx context.Context, client chat.ChatServiceClient){

	message := chat.Message{
		Body: "Hello people!",
		Lamport: safeL.Lamport,
	}

	response, err := client.SayHello(context.Background(), &message)
	

	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	safeL.CompareTimestamp(response.Lamport)

	log.Printf("Response from server: %s and %d", response.Body, safeL.Lamport)
}

func (l *SafeLamport) CompareTimestamp(otherClock int64) {
    l.mu.Lock()
    newValue := max(l.Lamport, otherClock)
    l.Lamport = newValue
    defer l.mu.Unlock()
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}