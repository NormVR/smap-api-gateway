package auth

import (
	"api-gateway/internal/models"
	"api-gateway/internal/models/auth"
	"context"
	"log"
	"time"

	authService "github.com/NormVR/smap_protobuf/gen/services/auth_service"
	userService "github.com/NormVR/smap_protobuf/gen/services/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	authConn   *grpc.ClientConn
	userConn   *grpc.ClientConn
	authClient authService.AuthServiceClient
	userClient userService.UserServiceClient
}

func New() *GrpcClient {
	// TODO: get address from config
	authConn, err := grpc.NewClient("auth-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	userConn, err := grpc.NewClient("user-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &GrpcClient{
		authConn:   authConn,
		userConn:   userConn,
		authClient: authService.NewAuthServiceClient(authConn),
		userClient: userService.NewUserServiceClient(userConn),
	}
}

func (c *GrpcClient) CreateUser(data *auth.RegisterData) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user, err := c.authClient.CreateUser(ctx, &authService.CreateUserRequest{
		Email:     data.Email,
		Username:  data.Username,
		Password:  data.Password,
		FirstName: data.Firstname,
		LastName:  data.Lastname,
	})

	if err != nil {
		return 0, err
	}

	return user.UserId, nil

}

func (c *GrpcClient) Login(data *auth.LoginData) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.authClient.Login(ctx, &authService.LoginRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return "", err
	}

	return response.JwtToken, nil
}

func (c *GrpcClient) ValidateToken(token string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.authClient.ValidateToken(ctx, &authService.TokenRequest{
		JwtToken: token,
	})

	if err != nil {
		return 0, err
	}

	return response.UserId, nil
}

func (c *GrpcClient) Logout(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := c.authClient.Logout(ctx, &authService.TokenRequest{
		JwtToken: token,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *GrpcClient) GetUser(id int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := c.userClient.GetUser(ctx, &userService.GetUserRequest{
		UserId: id,
	})

	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:        response.UserId,
		Email:     response.Email,
		Username:  response.Username,
		Firstname: response.Firstname,
		Lastname:  response.Lastname,
	}, nil
}

func (c *GrpcClient) Close() {
	if c.authConn != nil {
		c.authConn.Close()
	}

	if c.userConn != nil {
		c.userConn.Close()
	}
}
