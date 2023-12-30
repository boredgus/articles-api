package grpc

import (
	pb "user-management/grpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewArticleServiceClient(addr string) pb.ArticleServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("failed to connect to article service server: %v", err)
	}
	return Client{
		ArticleServiceClient: pb.NewArticleServiceClient(conn),
		conn:                 conn,
	}
}

type Client struct {
	pb.ArticleServiceClient
	conn *grpc.ClientConn
}

func (c *Client) CloseConn() error {
	return c.conn.Close()
}
