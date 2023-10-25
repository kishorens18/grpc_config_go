package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "go_config/proto"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) InsertData(ctx context.Context, req *pb.Request) (*emptypb.Empty, error) {

	var value interface{}
	err := json.Unmarshal(req.Value.Value, &value)
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("kishore").Collection("nithish")
	result, err := collection.InsertOne(context.Background(), bson.M{
		"key":   req.Key,
		"value": value,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(result)

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterMyServiceServer(s, &server{})
	if err2 := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err2)
	}
}
