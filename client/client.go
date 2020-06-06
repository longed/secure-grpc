package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	pb "secure-grpc/proto"
	"strconv"
	"time"
)

const (
	serverAddr = "127.0.0.1:4322"
)

func Request() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)

	if err != nil {
		log.Printf("client gRPC dial failed. %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := pb.NewCalculateClient(conn)

	stream, err := client.Division(ctx)
	if err != nil {
		log.Printf("client division stream failed. exit %v\n", err)
		return
	}
	defer stream.CloseSend()

	waitChannel := make(chan struct{})
	go sendAndRecv(stream, waitChannel)
	<-waitChannel
}

func sendAndRecv(stream pb.Calculate_DivisionClient, waitChannel chan struct{}) {
	for i := 0; i < 10; i += 1 {
		// send
		mapping := make(map[string]string)
		mapping["in"] = "ping " + strconv.Itoa(i)
		err := stream.Send(&pb.Req{Params: mapping})
		if err != nil {
			log.Fatal(fmt.Sprintf("stream client send failed. %v\n", err))
		}

		// receive
		rep, err := stream.Recv()
		if err == io.EOF {
			close(waitChannel)
			return
		}
		if err != nil {
			log.Fatal(fmt.Sprintf("both stream receive from server failed. %v\n", err))
		}
		log.Printf("received message from server: %v\n", rep.Result["out"])
	}

	stream.CloseSend()
}
