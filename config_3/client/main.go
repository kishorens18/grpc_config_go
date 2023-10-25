package main

import (
	"context"
	"encoding/json"
	pb "go_config/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

func main()  {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMyServiceClient(conn)

	key := "kishore"
	value := "hello world"

	valueJson, err := json.Marshal(value)
	if err != nil {
		log.Fatalf("Failed to marshal value to JSON: %v", err)
	}

	req := &pb.Request{
		Key: key,
		Value: &anypb.Any{
			Value: valueJson,
		},
	}

	_, err = client.InsertData(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to insert data: %v", err)
    }
}