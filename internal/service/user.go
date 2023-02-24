// Package service package with services
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/OVantsevich/User-Service/internal/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository repository interface for user service
//
//go:generate mockery --name=UserRepository --case=underscore --output=./mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	RefreshUser(ctx context.Context, id, token string) error
	DeleteUser(ctx context.Context, id string) error
}

// Expiration time of access token
const accessExp = time.Minute * 15

// Expiration time of refresh token
const refreshExp = time.Hour * 10

// Strength of password
const passwordStrength = 50

// User user service
type User struct {
	rps    UserRepository
	jwtKey []byte
}

// CustomClaims claims with login and role
type CustomClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// NewUserServiceClassic new user service
func NewUserServiceClassic(rps UserRepository, key string) *User {
	return &User{rps: rps, jwtKey: []byte(key)}
}

// Signup service signup
func (u *User) Signup(ctx context.Context, user *model.User) (accessToken, refreshToken string, userResult *model.User, err error) {
	if err = passwordvalidator.Validate(user.Password, passwordStrength); err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - Validate: %w", err)
	}
	user.Password, err = hashingPassword(user.Password)
	if err != nil {
		return "", "", nil, err
	}
	user.ID = uuid.New().String()
	userResult.ID = user.ID

	if userResult, err = u.rps.CreateUser(ctx, user); err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - CreateUser: %w", err)
	}

	accessToken, refreshToken, err = u.createJWT(ctx, userResult)
	if err != nil {
		return "", "", nil, fmt.Errorf("userService - Signup - createJWT: %w", err)
	}

	return
}

// Login service login
//
//nolint:dupl //just because
func (u *User) Login(ctx context.Context, login, password string) (accessToken, refreshToken string, err error) {
	var user *model.User

	if user, err = u.rps.GetUserByLogin(ctx, login); err != nil {
		return "", "", fmt.Errorf("userService - Login - GetUserByLogin: %w", err)
	}

	if !checkPasswordHash(user.Password, password) {
		return "", "", fmt.Errorf("userService - Login - Password invalid: %w", err)
	}

	accessToken, refreshToken, err = u.createJWT(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("userService - Login - createJWT: %w", err)
	}

	return
}

// Refresh service refresh
//
//nolint:dupl //just because
func (u *User) Refresh(ctx context.Context, id, userRefreshToken string) (accessToken, refreshToken string, err error) {
	var user *model.User

	if user, err = u.rps.GetUserByID(ctx, id); err != nil {
		return "", "", fmt.Errorf("userService - Refresh - GetUserByID: %w", err)
	}

	if user.Token != userRefreshToken {
		return "", "", fmt.Errorf("userService - Refresh - Token invalid: %w", err)
	}

	accessToken, refreshToken, err = u.createJWT(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("userService - Refresh - createJWT: %w", err)
	}

	return
}

// Update service update
func (u *User) Update(ctx context.Context, id string, user *model.User) (err error) {
	if err = u.rps.UpdateUser(ctx, id, user); err != nil {
		return fmt.Errorf("userService - Update - UpdateUser: %w", err)
	}

	return
}

// Delete service delete
func (u *User) Delete(ctx context.Context, id string) (err error) {
	if err = u.rps.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("userService - Delete - DeleteUser: %w", err)
	}

	return
}

// GetByLogin service get by login
func (u *User) GetByLogin(ctx context.Context, login string) (user *model.User, err error) {
	if user, err = u.rps.GetUserByLogin(ctx, login); err != nil {
		return nil, fmt.Errorf("userService - GetByLogin - Repository - GetByLogin: %w", err)
	}

	return
}

// GetByID service get by id
func (u *User) GetByID(ctx context.Context, id string) (user *model.User, err error) {
	if user, err = u.rps.GetUserByID(ctx, id); err != nil {
		return nil, fmt.Errorf("userService - GetByID - Repository - GetUserByID: %w", err)
	}

	return
}

func (u *User) createJWT(ctx context.Context, user *model.User) (accessTokenStr, refreshTokenStr string, err error) {
	accessClaims := &CustomClaims{
		user.ID,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExp)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err = accessToken.SignedString(u.jwtKey)
	if err != nil {
		return "", "", fmt.Errorf("userService - createJWT - SignedString: %w", err)
	}

	refreshClaims := &CustomClaims{
		user.ID,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExp)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err = refreshToken.SignedString(u.jwtKey)
	if err != nil {
		return "", "", fmt.Errorf("userService - createJWT - SignedString: %w", err)
	}

	err = u.rps.RefreshUser(ctx, user.ID, refreshTokenStr)
	if err != nil {
		return "", "", fmt.Errorf("userService - createJWT - RefreshUser: %w", err)
	}
	return accessTokenStr, refreshTokenStr, err
}

func hashingPassword(password string) (string, error) {
	hashedBytesPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("user - hashingPassword - GenerateFromPassword: %w", err)
	}
	return string(hashedBytesPassword), nil
}

func checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
