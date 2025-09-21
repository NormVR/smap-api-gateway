package services

import (
	"api-gateway/internal/grpc/auth"
	model "api-gateway/internal/models/auth"
	"fmt"
	"log"
)

type UserService struct {
	grpcClient *auth.GrpcClient
}

func NewUserService() *UserService {
	return &UserService{
		grpcClient: auth.New(),
	}
}

func (s *UserService) RegisterUser(userData *model.UserData) error {
	id, err := s.grpcClient.CreateUser(userData)
	defer s.grpcClient.Close()
	if err != nil {
		// TODO: process AlreadyExist scenario
		return fmt.Errorf("register error: %s", err)
	}

	log.Printf("Created user with id %d", id)
	return nil
}

func (s *UserService) LoginUser(loginData *model.LoginData) (map[string]string, error) {
	token, err := s.grpcClient.Login(loginData)
	defer s.grpcClient.Close()

	if err != nil {
		// TODO: process Wrong username or password scenario
		return nil, fmt.Errorf("login error: %s", err)
	}

	result := map[string]string{
		"token": token,
	}

	return result, nil
}
