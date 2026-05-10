package casbingrpc

import (
	"context"
	"fmt"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	authv1 "github.com/dungxbuif/RuntimeRoasters/runtime/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	serviceName string
	conn        *grpc.ClientConn
	authClient  authv1.AuthServiceClient
}

func NewAuthSnapshotClient(addr string, serviceName string) (casbin.AuthSnapshotClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &client{
		serviceName: serviceName,
		conn:        conn,
		authClient:  authv1.NewAuthServiceClient(conn),
	}, nil
}

func (c *client) GetFullSnapshot(ctx context.Context) ([]string, error) {
	resp, err := c.authClient.GetFullSnapshot(ctx, &authv1.GetFullSnapshotRequest{
		ServiceName: c.serviceName,
	})
	if err != nil {
		return nil, err
	}

	return resp.Policies, nil
}
