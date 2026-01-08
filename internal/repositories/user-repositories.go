package repositories

import (
	"context"
	"errors"
	"fmt"
	"go-auth-postgres/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

type UsersRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{
		DB: db,
	}
}

func (repo *UsersRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO go.users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var usuarioId uint
	now := time.Now()
	err := repo.DB.QueryRow(ctx, query, user.Username, user.Email, user.Password, now, now).Scan(&usuarioId)
	if err != nil {
		return fmt.Errorf("create, error while inserting user: %v", err)
	}

	return nil
}

func (repo *UsersRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	query := `SELECT * FROM go.users WHERE id = $1`
	err := repo.DB.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("FindByID, error while finding user: %v", err)
	}

	return &user, nil
}

func (repo *UsersRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, password FROM go.users WHERE email = $1`
	err := repo.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("FindByEmail, error while finding user: %v", err)
	}
	fmt.Println(user.Email)

	return &user, nil
}

func (repo *UsersRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	query := `UPDATE go.users SET username = $1, email = $2, password = $3, updated_at = $4 WHERE id = $5`
	_, err := repo.DB.Exec(ctx, query, user.Username, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("error while updating user: %v", err)
	}

	return nil
}
