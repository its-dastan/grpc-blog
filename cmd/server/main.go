package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/its-dastan/grpc-blog/pb"
	"github.com/its-dastan/grpc-blog/service"
	"google.golang.org/grpc"
)

const (
	secretKey     = "secret"
	tokenDuration = 5 * time.Minute
)

func runGRPCServer(authServer pb.AuthServiceServer, jwtManager *service.JWTManager, listener net.Listener) error {
	interceptor := service.NewAuthInterceptor(jwtManager)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	log.Printf("Start GRPC server on: %s", listener.Addr().String())
	return grpcServer.Serve(listener)
}

func runRESTServer(authServer pb.AuthServiceServer, jwtManager *service.JWTManager, listener net.Listener, endpoint string) error {
	fmt.Println("Hello REST")

	mux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, dialOptions)
	if err != nil {
		return err
	}

	log.Printf("Start REST server on: %s ", listener.Addr().String())
	return http.Serve(listener, mux)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	serverType := flag.String("type", "grpc", "type of server(grpc/rest)")
	endPoint := flag.String("endpoint", "", "grpc endpoint")
	flag.Parse()

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(jwtManager)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	if *serverType == "rest" {
		err = runRESTServer(authServer, jwtManager, listener, *endPoint)
	} else {
		err = runGRPCServer(authServer, jwtManager, listener)
	}

	if err != nil {
		log.Fatal("cannot start server")
	}
}
