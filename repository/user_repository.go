package repository

import (
	"context"
	"fmt"
	"time"

	"db/ent"
	"db/ent/user"

	"github.com/google/uuid"
)

// UserRepository представляет репозиторий для работы с пользователями
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository создаёт новый репозиторий пользователей
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// Create создаёт нового пользователя
func (r *UserRepository) Create(ctx context.Context, name string, email string, age int) (*ent.User, error) {
	u, err := r.client.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetAge(age).
		SetCreatedAt(time.Now()).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return u, nil
}

// CreateRandom создаёт пользователя со случайным email
func (r *UserRepository) CreateRandom(ctx context.Context) (*ent.User, error) {
	return r.Create(ctx, "John Doe", uuid.New().String(), 30)
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*ent.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed querying user by ID: %w", err)
	}
	return u, nil
}

// GetWithTransactions получает пользователя вместе с его транзакциями
func (r *UserRepository) GetWithTransactions(ctx context.Context, id int) (*ent.User, error) {
	userWithTx, err := r.client.User.
		Query().
		Where(user.ID(id)).
		WithTransactions().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user with transactions: %w", err)
	}

	return userWithTx, nil
}
