// Package middleware functions of middleware
package middleware

import (
	"context"
	"fmt"
	"github.com/OVantsevich/User-Service/internal/service"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JwtAuth checking token and attaching it to context
func JwtAuth(keyFunc func(token *jwt.Token) (interface{}, error)) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod != "/UserService/Signup" && info.FullMethod != "/UserService/Login" && info.FullMethod != "/UserService/Refresh" {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
			}

			authHeader, ok := md["authorization"]
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
			}

			token := authHeader[0]

			claims, err := verify(token, keyFunc)
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, "user", claims)
		}

		h, err := handler(ctx, req)

		return h, err
	})
}

func verify(token string, keyFunc func(token *jwt.Token) (interface{}, error)) (claims *service.CustomClaims, err error) {
	claims = &service.CustomClaims{}

	_, err = jwt.ParseWithClaims(
		token,
		claims,
		keyFunc,
	)
	if err != nil {
		err = fmt.Errorf("invalid token: %w", err)
	}

	return
}
