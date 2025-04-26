package service

import (
	"context"
	"accounting/ent"
	"accounting/repository"
	"fmt"
)

// UserService represents a service for working with users
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(client *ent.Client) *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(client),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, name, email string, age int) (*ent.User, error) {
	user, err := s.userRepo.Create(ctx, name, email, age)
	if err != nil {
		return nil, fmt.Errorf("user service - create user: %w", err)
	}
	return user, nil
}

// CreateRandomUser creates a user with random data (for testing)
func (s *UserService) CreateRandomUser(ctx context.Context) (*ent.User, error) {
	user, err := s.userRepo.CreateRandom(ctx)
	if err != nil {
		return nil, fmt.Errorf("user service - create random user: %w", err)
	}
	return user, nil
}

// GetUserByID gets a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*ent.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user service - get user by id: %w", err)
	}
	return user, nil
}

// GetUserWithTransactions gets a user together with their transactions
func (s *UserService) GetUserWithTransactions(ctx context.Context, id int) (*ent.User, error) {
	user, err := s.userRepo.GetWithTransactions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user service - get user with transactions: %w", err)
	}
	return user, nil
}
