package services

import (
	"api-gateway/internal/grpc/auth"
	userModel "api-gateway/internal/models"
	model "api-gateway/internal/models/auth"
	"api-gateway/internal/repository/redis"
	"log"

	"github.com/google/uuid"
)

type UserService struct {
	grpcClient *auth.GrpcClient
	redisRepo  *redis.RedisRepository
}

func NewUserService(grpcClient *auth.GrpcClient, redisRepo *redis.RedisRepository) *UserService {
	return &UserService{
		grpcClient: grpcClient,
		redisRepo:  redisRepo,
	}
}

func (s *UserService) RegisterUser(userData *model.AuthData) error {
	id, err := s.grpcClient.CreateUser(userData)
	if err != nil {
		return err
	}

	log.Printf("Created user with id %d", id)
	return nil
}

func (s *UserService) LoginUser(loginData *model.AuthData) (map[string]string, error) {
	token, err := s.grpcClient.Login(loginData)

	if err != nil {
		return nil, err
	}

	result := map[string]string{
		"token": token,
	}

	return result, nil
}

func (s *UserService) Logout(tokenString string) error {
	err := s.grpcClient.Logout(tokenString)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *UserService) ValidateToken(token string) (uuid.UUID, error) {
	userId, err := s.redisRepo.GetUserId("token:" + token)

	if err != nil {
		log.Printf(err.Error())
	}

	if userId == uuid.Nil {
		userId, err = s.grpcClient.ValidateToken(token)
		if err != nil {
			return uuid.Nil, err
		}

		if userId == uuid.Nil {
			return uuid.Nil, nil
		}
	}

	return userId, nil
}

func (s *UserService) GetUser(userId uuid.UUID) (*userModel.User, error) {
	user, err := s.grpcClient.GetUser(userId)

	if err != nil {
		return nil, err
	}

	return user, nil
}
