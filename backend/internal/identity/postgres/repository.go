package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"support-go/backend/internal/identity"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) CreateUser(ctx context.Context, user identity.User) (identity.User, error) {
	query := `
		INSERT INTO users (id, email, password_hash, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, password_hash, role, status, created_at, updated_at
	`

	var created identity.User
	err := repository.db.QueryRow(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&created.ID,
		&created.Email,
		&created.PasswordHash,
		&created.Role,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return identity.User{}, identity.ErrEmailAlreadyExists
		}
		return identity.User{}, err
	}

	return created, nil
}

func (repository *Repository) GetUserByEmail(ctx context.Context, email string) (identity.User, error) {
	return repository.getOne(ctx, `
		SELECT id, email, password_hash, role, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)
}

func (repository *Repository) GetUserByID(ctx context.Context, id string) (identity.User, error) {
	return repository.getOne(ctx, `
		SELECT id, email, password_hash, role, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)
}

func (repository *Repository) getOne(ctx context.Context, query string, value string) (identity.User, error) {
	var user identity.User
	err := repository.db.QueryRow(ctx, query, value).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return identity.User{}, identity.ErrUserNotFound
	}
	if err != nil {
		return identity.User{}, err
	}

	return user, nil
}
