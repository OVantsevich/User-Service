// Package repository User
package repository

import (
	"context"
	"fmt"
	"time"

	"User-Service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// userID id
const userID = "user"

// User postgres entity
type User struct {
	Pool *pgxpool.Pool
}

// NewUser creating new User repository
func NewUser(pool *pgxpool.Pool) *User {
	return &User{Pool: pool}
}

// CreateUser create user
func (r *User) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	user.Created = time.Now()
	user.Updated = time.Now()
	_, err := r.Pool.Exec(ctx,
		"insert into users (id, login, email, role, password, name, age) values ($1, $2, $3, $4, $5, $6, $7) returning role;",
		user.ID, user.Login, user.Email, userID, user.Password, user.Name, user.Age)
	if err != nil {
		return nil, fmt.Errorf("user - CreateUser - Exec: %w", err)
	}

	return user, nil
}

// GetUserByLogin get user by login
func (r *User) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	user := model.User{}
	err := r.Pool.QueryRow(ctx, `select u.id, u.name, u.age, u.login, u.password, u.token, u.email, r.name
									from users u
											 join roles r on r.id = u.role
									where u.login = $1 and u.deleted=false`, login).Scan(
		&user.ID, &user.Name, &user.Age, &user.Login, &user.Password, &user.Token, &user.Email, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("user - GetUserByLogin - QueryRow: %w", err)
	}

	return &user, nil
}

// GetUserByID get user by login
func (r *User) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := model.User{}
	err := r.Pool.QueryRow(ctx, `select u.id, u.name, u.age, u.login, u.password, u.token, u.email, r.name
									from users u
											 join roles r on r.id = u.role
									where u.id = $1 and u.deleted=false`, id).Scan(
		&user.ID, &user.Name, &user.Age, &user.Login, &user.Password, &user.Token, &user.Email, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("user - GetUserByID - QueryRow: %w", err)
	}

	return &user, nil
}

// UpdateUser update user
func (r *User) UpdateUser(ctx context.Context, id string, user *model.User) error {
	var idCheck int
	err := r.Pool.QueryRow(ctx, "update users set email=$1, name=$2, age=$3, updated=$4 where id=$5 and deleted=false returning id",
		user.Email, user.Name, user.Age, user.Updated, id).Scan(&idCheck)
	if err != nil {
		return fmt.Errorf("user - UpdateUser - Exec: %w", err)
	}

	return nil
}

// RefreshUser refresh user
func (r *User) RefreshUser(ctx context.Context, id, token string) error {
	var idCheck int
	err := r.Pool.QueryRow(ctx, "update users set token=$1, updated=$2 where id=$3 and deleted=false returning id",
		token, time.Now(), id).Scan(&idCheck)
	if err != nil {
		return fmt.Errorf("user - RefreshUser - Exec: %w", err)
	}

	return nil
}

// DeleteUser delete user
func (r *User) DeleteUser(ctx context.Context, id string) error {
	var idCheck int
	err := r.Pool.QueryRow(ctx, "update users set Deleted=true, updated=$1 where id=$2 and deleted=false returning id",
		time.Now(), id).Scan(&idCheck)
	if err != nil {
		return fmt.Errorf("user - DeleteUser - Exec: %w", err)
	}

	return nil
}
