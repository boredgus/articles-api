package main

import (
	"flag"
	"fmt"
	"net"

	"user-management/config"
	pb "user-management/grpc"
	server "user-management/grpc/server"
	"user-management/internal/gateways"
	"user-management/internal/models"
	"user-management/pkg/db"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func init() {
	config.InitConfig()
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterArticleServiceServer(s, server.NewArticleService(models.NewArticleModel(
		gateways.NewCachedArticleRepository(
			gateways.NewArticleRepository(
				db.NewMySQLStore(),
				db.NewClickHouseStore()),
			db.NewRedisStore(),
		))))

	logrus.Printf("article service server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
