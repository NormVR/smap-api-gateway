package services

import (
	"api-gateway/internal/grpc/auth"
	model "api-gateway/internal/models/auth"
	"api-gateway/internal/repository/redis"
	"log"
)

type UserService struct {
	grpcClient *auth.GrpcClient
	redisRepo  *redis.RedisRepository
}

func NewUserService() *UserService {
	return &UserService{
		grpcClient: auth.New(),
		redisRepo:  redis.NewRedisRepository(),
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

func (s *UserService) ValidateToken(token string) (int64, error) {
	userId, err := s.redisRepo.GetUserId(token)

	if err != nil {
		log.Printf(err.Error())
	}

	if userId == 0 {
		userId, err = s.grpcClient.ValidateToken(token)
		if err != nil {
			return 0, err
		}

		if userId == 0 {
			return 0, nil
		}
	}

	return userId, nil
}
