package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"codebase/internal/bob"
	"codebase/internal/models"
	"codebase/internal/service/utils"
	"codebase/pkg/errorx"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do"
	"golang.org/x/crypto/bcrypt"
)

type ServiceAuth interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
}

type serviceAuth struct {
	Utils     *utils.ServiceUtils
	Datastore models.Datastore
}

func NewServiceAuth(container *do.Injector) (ServiceAuth, error) {
	utilsService, err := do.Invoke[*utils.ServiceUtils](container)
	if err != nil {
		return nil, err
	}

	datastore, err := do.Invoke[models.Datastore](container)
	if err != nil {
		return nil, err
	}

	return &serviceAuth{
		Utils:     utilsService,
		Datastore: datastore,
	}, nil
}

func (s *serviceAuth) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	existing, _ := s.Datastore.FindUserByUsername(ctx, req.Username)
	if existing != nil {
		return nil, errors.New("username already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New()
	param := &bob.UserSetter{
		ID:        omit.From(userID),
		Username:  omit.From(req.Username),
		Password:  omit.From(string(hashed)),
		Email:     omit.From(req.Email),
		FirstName: omitnull.From(req.FirstName),
		LastName:  omitnull.From(req.LastName),
		IsActive:  omit.From(true), // auto active for simplicity
		CreatedAt: omit.From(time.Now()),
		UpdatedAt: omit.From(time.Now()),
	}

	user, err := s.Datastore.CreateUser(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *serviceAuth) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.Datastore.FindRawUserByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return nil, errorx.Wrap(fmt.Errorf("invalid username or password"), errorx.NotExist)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errorx.Wrap(fmt.Errorf("invalid username or password"), errorx.Invalid)
	}

	claims := jwt.MapClaims{
		"id":         user.ID.String(),
		"username":   user.Username,
		"email":      user.Email,
		"first_name": user.FirstName.GetOrZero(),
		"last_name":  user.LastName.GetOrZero(),
		"is_active":  user.IsActive,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "fallback_secret_codebase"
	}
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &models.AuthResponse{
		Token: t,
		User:  models.UserBobToRaw(user),
	}, nil
}
