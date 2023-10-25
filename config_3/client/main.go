package main

import (
	"context"
	pb "go_config/proto"
	"log"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func main() {
	conc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer conc.Close()

	client := pb.NewMyServiceClient(conc)

	postmanData := map[string]interface{}{
		"Key":   "hello world",
		"Value": true,
	}

	dataBytes, err := proto.Marshal(postmanData)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	req := &pb.InsertRequest{
		Data: &any.Any{
			Value: dataBytes,
		},
	}

	res, err := client.InsertData(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to insert data via gRPC: %v", err)
	}

	if res.Success {
		log.Printf("Data inserted successfully via gRPC: %s", res.Message)
	} else {
		log.Printf("Failed to insert data via gRPC: %s", res.Message)
	}
}
