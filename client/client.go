package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
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
	certificate, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("append err!")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs: certPool,
		ServerName: "server.io",
	})

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))

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
