package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"secure-grpc/client"
	pb "secure-grpc/proto"
	"secure-grpc/server"
	"time"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("append err")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs: certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	})

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
