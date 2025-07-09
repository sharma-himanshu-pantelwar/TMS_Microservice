package sessionclient

import (
	"context"
	"fmt"
	pb "task_service/src/internal/interfaces/input/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.SessionValidatorClient
}

func NewClient(userServiceAddress string) (*Client, error) {
	conn, err := grpc.Dial(userServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewSessionValidatorClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}
func (c *Client) ValidateSession(ctx context.Context, sessionID string) (bool, string, error) {
	resp, err := c.client.ValidateSession(ctx, &pb.ValidateSessionRequest{SessionId: sessionID})
	fmt.Println("Response in grpcClient: ", resp)
	if err != nil {
		return false, "", err
	}
	return resp.Valid, resp.UserId, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
