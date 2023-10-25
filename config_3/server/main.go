package main

import (
	"context"
	pb "go_config/proto"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
)

type server struct {
	pb.UnimplementedMyServiceServer
}

var (
	ConfigCollection *mongo.Collection
)

func Connectdatabase() (*mongo.Client, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancle()
	mongoConnection := options.Client().ApplyURI("mongodb://localhost:27017")
	MongoClient, err := mongo.Connect(ctx, mongoConnection)
	if err != nil {
		log.Fatal()
		return nil, err
	}

	err1 := MongoClient.Ping(ctx, readpref.Primary())
	if err1 != nil {
		return nil, err1
	}
	return MongoClient, nil

}

func GetCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}

func (s *server) InserData(ctx context.Context, req *pb.InsertRequest) (*pb.InsertResponse, error) {
	data := req.GetData()


	// var unknownData pb.UnknownData
	// if err := anypb.UnmarshalTo(data, &unknownData, proto.UnmarshalOptions{}); err != nil {
	// 	return nil, err
	// }


	var originalData map[string]interface{}
	err := proto.Unmarshal(data.Value, &originalData)
	if err != nil {
		return nil, err
	}
	_, err2 := ConfigCollection.InsertOne(context.TODO(), originalData)
	if err2 != nil {
		return nil, err2
	}
	return &pb.InsertResponse{
		Success: true,
		Message: "Data inserted succesfully",
	}, nil
}

func main() {
	mongoClient, err := Connectdatabase()
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(context.TODO())

	ConfigCollection = GetCollection(mongoClient, "GoConfig", "config")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterMyServiceServer(s, &server{})
	log.Println("gRPC server listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
