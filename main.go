package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"secure-grpc/client"
	pb "secure-grpc/proto"
	"secure-grpc/server"
	"time"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 4322))
	if err != nil {
		log.Printf("listen port failed. %v\n", err)
		return
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterCalculateServer(grpcServer, &server.CalculatorServer{})

	go func() {
		time.Sleep(time.Second * 2)
		client.Request()
	}()

	grpcServer.Serve(lis)
}
