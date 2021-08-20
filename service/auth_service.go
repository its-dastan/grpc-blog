package service

import "github.com/its-dastan/grpc-blog/pb"

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
}
