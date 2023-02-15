// Package main main
package main

import (
	"context"
	"fmt"
	"net"

	"User-Service/internal/config"
	"User-Service/internal/handler"
	"User-Service/internal/middleware"
	"User-Service/internal/repository"
	"User-Service/internal/service"
	pr "User-Service/proto"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.NewMainConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	listen, err := net.Listen("tcp", "localhost:12344")
	if err != nil {
		defer logrus.Fatalf("error while listening port: %e", err)
	}
	fmt.Println("Server successfully started on port :12344...")

	var repos service.UserRepository
	repos, err = dbConnection(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	defer closePool(repos)

	userService := service.NewUserServiceClassic(repos, cfg.JwtKey)

	ns := grpc.NewServer(middleware.JwtAuth(func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtKey), nil
	}))
	server := handler.NewUserHandlerClassic(userService, cfg.JwtKey)
	pr.RegisterUserServiceServer(ns, server)

	if err = ns.Serve(listen); err != nil {
		defer logrus.Fatalf("error while listening server: %e", err)
	}
}

func dbConnection(cfg *config.MainConfig) (service.UserRepository, error) {
	pgURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword,
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)

	pool, err := pgxpool.New(context.Background(), pgURL)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration data: %v", err)
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database not responding: %v", err)
	}
	return repository.NewUser(pool), nil
}

func closePool(r interface{}) {
	p := r.(repository.User)
	if p.Pool != nil {
		p.Pool.Close()
	}
}
