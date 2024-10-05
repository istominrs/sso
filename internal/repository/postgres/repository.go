package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"sso/internal/domain/models"
	"sso/internal/repository"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(dsn string) (*Repository, error) {
	const op = "repository.postgresql.New"

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{pool: pool}, nil
}

func (r *Repository) Close() {
	r.pool.Close()
}

func (r *Repository) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "repository.postgresql.SaveUser"

	row := r.pool.QueryRow(ctx, "INSERT INTO users(email, pass_hash) VALUES ($1, $2) RETURNING id", email, passHash)

	var id int64
	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repository) User(ctx context.Context, email string) (models.User, error) {
	const op = "repository.postgresql.User"

	row := r.pool.QueryRow(ctx, "SELECT id, email, pass_hash FROM users WHERE email = $1", email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *Repository) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "repository.postgresql.IsAdmin"

	row := r.pool.QueryRow(ctx, "SELECT COUNT(1) FROM admins WHERE user_id = $1", userID)

	var count int
	if err := row.Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count == 1, nil
}

func (r *Repository) App(ctx context.Context, id int) (models.App, error) {
	const op = "repository.postgresql.App"

	row := r.pool.QueryRow(ctx, "SELECT id, name, secret FROM apps WHERE id = $1", id)

	var app models.App
	if err := row.Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, repository.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
