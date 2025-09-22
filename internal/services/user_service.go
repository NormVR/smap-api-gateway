package services

import (
	"api-gateway/internal/grpc/auth"
	model "api-gateway/internal/models/auth"
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
	if err != nil {
		return err
	}

	log.Printf("Created user with id %d", id)
	return nil
}

func (s *UserService) LoginUser(loginData *model.LoginData) (map[string]string, error) {
	token, err := s.grpcClient.Login(loginData)

	if err != nil {
		return nil, err
	}

	result := map[string]string{
		"token": token,
	}

	return result, nil
}
