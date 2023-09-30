package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jdiuwe/api-server-demo/types"
	_ "github.com/lib/pq" //nolint: revive
	"github.com/rs/zerolog"
)

const (
	connStr        = "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"
	selectUserStmt = "SELECT id, name FROM users WHERE id = $1"
	insertUserStmt = "INSERT INTO users (id, name) VALUES ($1, $2)"
)

var (
	ErrGetUserQuery       = errors.New("error getting user")
	ErrUserNotFound       = errors.New("user not found")
	ErrCouldNotCreateUser = errors.New("could not create user")
)

type Store interface {
	GetUser(ctx context.Context, id string) (*types.User, error)
	CreateUser(ctx context.Context, user *types.User) error
}

type PostgresStore struct {
	*sql.DB
	logger zerolog.Logger
}

func NewPostgresStore(logger zerolog.Logger, cfg *Config) (*PostgresStore, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("db configuration error: %w", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf(connStr, cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password))
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}

	return &PostgresStore{db, logger}, nil
}

func (s *PostgresStore) GetUser(ctx context.Context, id string) (*types.User, error) {
	stmt, err := s.PrepareContext(ctx, selectUserStmt)
	if err != nil {
		s.logger.Err(err).Msg("PrepareContext error")

		return nil, ErrGetUserQuery
	}
	defer stmt.Close()

	var user types.User

	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		s.logger.Err(err).Msg("ExecContext error")

		return &user, ErrGetUserQuery
	}

	return &user, nil
}

func (s *PostgresStore) CreateUser(ctx context.Context, user *types.User) error {
	stmt, err := s.PrepareContext(ctx, insertUserStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.ID, user.Name)
	if err != nil {
		s.logger.Err(err).Msg("ExecContext error")

		return ErrCouldNotCreateUser
	}

	s.logger.Debug().Msgf("PostgresStore.CreateUser: %+v", user)

	return nil
}
