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

// func (c *Client) ValidateSession(ctx context.Context, sessionID string) (bool, string, error) {
// 	resp, err := c.client.ValidateSession(ctx, &pb.ValidateSessionRequest{SessionId: sessionID})
// 	fmt.Println("Response in grpcClient: ", resp)
// 	if err != nil {
// 		fmt.Println("err is ", err)
// 		return false, "", err
// 	}
// 	return resp.Valid, resp.UserId, nil
// }

func (c *Client) ValidateSession(ctx context.Context, sessionID string) (bool, int64, error) {
	resp, err := c.client.ValidateSession(ctx, &pb.ValidateSessionRequest{SessionId: sessionID})
	fmt.Println("Response in grpcClient: ", resp)
	if err != nil {
		fmt.Println("err is ", err)
		return false, 0, err
	}
	if resp == nil {
		fmt.Println("nil response from ValidateSession")
		return false, 0, fmt.Errorf("recieved nil response from ValidateSession")
	}
	return resp.Valid, resp.UserId, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
