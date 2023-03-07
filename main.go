// Package main main
package main

import (
	"context"
	"fmt"
	"net"

	"github.com/OVantsevich/User-Service/internal/config"
	"github.com/OVantsevich/User-Service/internal/handler"
	"github.com/OVantsevich/User-Service/internal/repository"
	"github.com/OVantsevich/User-Service/internal/service"
	pr "github.com/OVantsevich/User-Service/proto"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.NewMainConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		logrus.Fatalf("error while listening port: %v", err)
	}

	var repos service.UserRepository
	repos, err = dbConnection(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	defer closePool(repos)

	userService := service.NewUserServiceClassic(repos, cfg.JwtKey)

	//ns := grpc.NewServer(middleware.JwtAuth(func(token *jwt.Token) (interface{}, error) {
	//	return []byte(cfg.JwtKey), nil
	//}))
	ns := grpc.NewServer()
	server := handler.NewUserHandlerClassic(userService, cfg.JwtKey)
	pr.RegisterUserServiceServer(ns, server)

	if err = ns.Serve(listen); err != nil {
		logrus.Fatalf("error while listening server: %v", err)
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
