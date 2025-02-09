package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"sso/sso/internal/domain/models"
	"sso/sso/internal/lib/jwt"
	"sso/sso/internal/storage"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user exists")
	ErrInvalidAppID       = errors.New("invalid app id")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "internal.services.auth.Login()"

	a.log.Info(
		"attempting to login user",
		slog.String("op", op),
	)

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn(
				storage.ErrUserNotFound.Error()+": "+err.Error(),
				slog.String("op", op),
			)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error(
			"failed to get user: "+err.Error(),
			slog.String("op", op),
		)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		a.log.Warn(
			ErrInvalidCredentials.Error()+": "+err.Error(),
			slog.String("op", op),
		)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		a.log.Error(
			"failed to get app id: "+err.Error(),
			slog.String("op", op),
		)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info(
		"user logged in",
		slog.String("op", op),
	)

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error(
			"failed to generate token: "+err.Error(),
			slog.String("op", op),
		)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "internal.services.auth.RegisterNewUser()"

	a.log.Info(
		"registering user",
		slog.String("op", op),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error(
			"failed to generate password hash: "+err.Error(),
			slog.String("op", op),
		)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn(
				storage.ErrUserNotFound.Error()+": "+err.Error(),
				slog.String("op", op),
			)

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		a.log.Error(
			"failed to save user: "+err.Error(),
			slog.String("op", op),
		)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info(
		"user registered",
		slog.String("op", op),
	)

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "internal.services.auth.IsAdmin()"

	a.log.Info(
		"checking if user is admin",
		slog.String("op", op),
	)

	ok, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn(
				storage.ErrUserNotFound.Error()+": "+err.Error(),
				slog.String("op", op),
			)

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		a.log.Error(
			"failed to check if user is admin: "+err.Error(),
			slog.String("op", op),
		)

		return false, fmt.Errorf("%s: %w", op, err)
	}

	if !ok {
		a.log.Warn(
			"user is not admin",
			slog.String("op", op),
		)

		return false, nil
	}

	a.log.Info(
		"user is admin",
		slog.String("op", op),
	)

	return true, nil
}
