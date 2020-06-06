package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"secure-grpc/client"
	pb "secure-grpc/proto"
	"secure-grpc/server"
	"time"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 4322))
	if err != nil {
		log.Printf("listen port failed. %v\n", err)
		return
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCalculateServer(grpcServer, &server.CalculatorServer{})

	go func() {
		time.Sleep(time.Second * 2)
		client.Request()
	}()

	grpcServer.Serve(lis)
}
