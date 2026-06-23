package service

import (
	"context"
	"errors"
	"time"

	"github.com/exizzxc/barberguide-backend/internal/dto"
	"github.com/exizzxc/barberguide-backend/internal/model"
	"github.com/exizzxc/barberguide-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GetMe(ctx context.Context, userID uint) (*dto.UserDTO, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Проверяем что email ещё не занят
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// Хешируем пароль — в базе никогда не хранится открытый пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаём пользователя
	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Role:         model.RoleUser, // роль всегда user при регистрации
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Генерируем JWT-токен
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: token,
		User: dto.UserDTO{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  string(user.Role),
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Ищем пользователя по email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Генерируем JWT-токен
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: token,
		User: dto.UserDTO{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  string(user.Role),
		},
	}, nil
}

func (s *authService) GetMe(ctx context.Context, userID uint) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserDTO{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
	}, nil
}

func (s *authService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    string(user.Role),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
