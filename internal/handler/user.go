// Package handler handler
package handler

import (
	"context"
	"fmt"
	"github.com/OVantsevich/User-Service/internal/model"
	"github.com/OVantsevich/User-Service/internal/service"
	pr "github.com/OVantsevich/User-Service/proto"

	"github.com/sirupsen/logrus"
)

// UserService service interface for user handler
//
//go:generate mockery --name=UserService --case=underscore --output=./mocks
type UserService interface {
	Signup(ctx context.Context, user *model.User) (string, string, *model.User, error)
	Login(ctx context.Context, login, password string) (string, string, error)
	Refresh(ctx context.Context, id, userRefreshToken string) (string, string, error)
	Update(ctx context.Context, id string, user *model.User) error
	Delete(ctx context.Context, id string) error

	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}

// User handler
type User struct {
	pr.UnimplementedUserServiceServer
	service UserService
	jwtKey  string
}

// NewUserHandlerClassic new user handler
func NewUserHandlerClassic(s UserService, key string) *User {
	return &User{service: s, jwtKey: key}
}

// Signup handler signup
func (h *User) Signup(ctx context.Context, request *pr.SignupRequest) (response *pr.SignupResponse, err error) {
	user := &model.User{
		Login:    request.Login,
		Email:    request.Email,
		Password: request.Password,
		Name:     request.Name,
		Age:      int(request.Age),
	}

	var userResponse *model.User
	response = &pr.SignupResponse{}
	response.AccessToken, response.RefreshToken, userResponse, err = h.service.Signup(ctx, user)
	if err != nil {
		err = fmt.Errorf("userHandler - Signup - Signup: %w", err)
		logrus.Error(err)
		return
	}
	response.User = &pr.User{
		Id:    userResponse.ID,
		Login: userResponse.Login,
		Email: userResponse.Email,
		Name:  userResponse.Name,
		Age:   int32(userResponse.Age),
	}

	return
}

// Login handler login
func (h *User) Login(ctx context.Context, request *pr.LoginRequest) (response *pr.LoginResponse, err error) {
	response = &pr.LoginResponse{}
	response.AccessToken, response.RefreshToken, err = h.service.Login(ctx, request.Login, request.Password)
	if err != nil {
		err = fmt.Errorf("userHandler - Login - Login: %w", err)
		logrus.Error(err)
		return
	}

	return
}

// Refresh handler refresh
func (h *User) Refresh(ctx context.Context, request *pr.RefreshRequest) (response *pr.RefreshResponse, err error) {
	response = &pr.RefreshResponse{}
	response.AccessToken, response.RefreshToken, err = h.service.Refresh(ctx, request.Id, request.RefreshToken)
	if err != nil {
		err = fmt.Errorf("userHandler - Refresh - Refresh: %w", err)
		logrus.Error(err)
		return
	}

	return
}

// Update handler update
func (h *User) Update(ctx context.Context, request *pr.UpdateRequest) (response *pr.UpdateResponse, err error) {
	var claims = ctx.Value("user").(*service.CustomClaims)

	user := &model.User{
		Email: request.Email,
		Name:  request.Name,
		Age:   int(request.Age),
	}
	response = &pr.UpdateResponse{}
	err = h.service.Update(ctx, claims.ID, user)
	if err != nil {
		err = fmt.Errorf("userHandler - Update - Update: %w", err)
		logrus.Error(err)
		return
	}
	response.Success = true

	return
}

// Delete handler delete
func (h *User) Delete(ctx context.Context, _ *pr.Request) (response *pr.DeleteResponse, err error) {
	var claims = ctx.Value("user").(*service.CustomClaims)

	response = &pr.DeleteResponse{}
	err = h.service.Delete(ctx, claims.ID)
	if err != nil {
		err = fmt.Errorf("userHandler - Delete - Delete: %w", err)
		logrus.Error(err)
		return
	}
	response.Success = true

	return
}

// UserByID handler user by login
func (h *User) UserByID(ctx context.Context, request *pr.UserByIdRequest) (response *pr.UserByIdResponse, err error) {
	var claims = ctx.Value("user").(*service.CustomClaims)

	if claims.Role != "admin" {
		err = fmt.Errorf("access denied")
		logrus.Error(err)
		return
	}

	response = &pr.UserByIdResponse{}
	var user *model.User
	user, err = h.service.GetByID(ctx, request.ID)
	if err != nil {
		err = fmt.Errorf("userHandler - UserByID - GetByID: %w", err)
		logrus.Error(err)
		return
	}
	response.User = &pr.User{
		Id:    user.ID,
		Login: user.Login,
		Email: user.Email,
		Name:  user.Name,
		Age:   int32(user.Age),
	}

	return
}
