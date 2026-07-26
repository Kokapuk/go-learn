package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

var errUsernameTaken = errors.New("Username is already taken")

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) createUser(
	ctx context.Context,
	req createUserRequest,
) (User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return User{}, err
	}

	params := createUserParams{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
	}

	user, err := s.repository.createUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_key" {
			return User{}, errUsernameTaken
		}

		return User{}, err
	}

	return user, nil
}
