package server

import (
	"fmt"
	"io"
	"log"
	pb "secure-grpc/proto"
)

type CalculatorServer struct {
}

func (s *CalculatorServer) Division(stream pb.Calculate_DivisionServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("server received end event. exit")
				return nil
			}
			fmt.Printf("exix for recv err. %v ", err)
			return err
		}

		log.Println(in.Params["in"])
		result := make(map[string]string)
		result["out"] = "pong " + in.Params["in"]
		if err = stream.Send(&pb.Rep{Result: result}); err != nil {
			fmt.Printf("server send err. %v\n", err)
		}
	}
}
