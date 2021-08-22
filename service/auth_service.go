package service

import (
	"context"
	"github.com/its-dastan/grpc-blog/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	JWTManager *JWTManager
}

func NewAuthServer(jwtManager *JWTManager) *AuthServer {
	return &AuthServer{JWTManager: jwtManager}
}

func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user := &User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	token, err := user.Login(server)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	res := &pb.AuthResponse{AccessToken: string(token)}
	return res, nil
}

func (server *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user := &User{
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		MobileNumber: req.GetMobileNumber(),
		Password:     req.GetPassword(),
	}

	token, err := user.Register(server)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	res := &pb.AuthResponse{AccessToken: string(token)}
	return res, nil
}
