package grpc_clients

import (
	pb "go-clinet-locations/shared/proto/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

type locationServiceClient struct {
	Client pb.LocationServiceClient
	conn   *grpc.ClientConn
}

func NewLocationServiceClient() (*locationServiceClient, error) {
	locationServiceURL := os.Getenv("LOCATION_SERVICE_URL")
	if locationServiceURL == "" {
		locationServiceURL = "location-history-service:9092"
	}
	conn, err := grpc.NewClient(locationServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewLocationServiceClient(conn)

	return &locationServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *locationServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
