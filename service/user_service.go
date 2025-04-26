package service

import (
	"context"
	"db/ent"
	"db/repository"
	"fmt"
)

// UserService представляет сервис для работы с пользователями
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService создает новый сервис пользователей
func NewUserService(client *ent.Client) *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(client),
	}
}

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(ctx context.Context, name, email string, age int) (*ent.User, error) {
	user, err := s.userRepo.Create(ctx, name, email, age)
	if err != nil {
		return nil, fmt.Errorf("user service - create user: %w", err)
	}
	return user, nil
}

// CreateRandomUser создаёт пользователя со случайными данными (для тестирования)
func (s *UserService) CreateRandomUser(ctx context.Context) (*ent.User, error) {
	user, err := s.userRepo.CreateRandom(ctx)
	if err != nil {
		return nil, fmt.Errorf("user service - create random user: %w", err)
	}
	return user, nil
}

// GetUserByID получает пользователя по его ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*ent.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user service - get user by id: %w", err)
	}
	return user, nil
}

// GetUserWithTransactions получает пользователя вместе с его транзакциями
func (s *UserService) GetUserWithTransactions(ctx context.Context, id int) (*ent.User, error) {
	user, err := s.userRepo.GetWithTransactions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user service - get user with transactions: %w", err)
	}
	return user, nil
}
