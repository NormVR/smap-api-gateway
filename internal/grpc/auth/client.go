package auth

import (
	"api-gateway/internal/models/auth"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	user_service "github.com/NormVR/smap_protobuf/gen"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	client user_service.AuthServiceClient
}

func New() *GrpcClient {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &GrpcClient{
		conn:   conn,
		client: user_service.NewAuthServiceClient(conn),
	}
}

func (c *GrpcClient) CreateUser(data *auth.UserData) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, err := c.client.CreateUser(ctx, &user_service.CreateUserRequest{
		Email:     data.Email,
		Username:  data.Username,
		Password:  data.Password,
		FirstName: data.Firstname,
		LastName:  data.Lastname,
	})

	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return 0, err
	}

	log.Println(user.UserId)
	return user.UserId, nil

}

func (c *GrpcClient) Login(data *auth.LoginData) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.client.Login(ctx, &user_service.LoginRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return "", err
	}
	return response.JwtToken, nil
}

func (c *GrpcClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
