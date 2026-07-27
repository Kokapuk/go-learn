package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

var errUsernameTaken = errors.New("Username is already taken")
var errInvalidUsernameOrPassword = errors.New("Username or password is invalid")

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) signUp(
	ctx context.Context,
	req signUpRequest,
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

func (s *Service) signIn(ctx context.Context, req signInRequest) (User, error) {
	user, err := s.repository.getUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, errInvalidUsernameOrPassword
		}

		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return User{}, errInvalidUsernameOrPassword
	}

	return User{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt}, nil
}

func (s *Service) getSelf(ctx context.Context, userID string) (User, error) {
	return s.repository.getPublicUserByID(ctx, userID)
}
