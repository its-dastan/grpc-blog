gen: 
	protoc -I=proto --go_out=. --go-grpc_out=. proto/*.proto --grpc-gateway_out=. --swagger_out=:swagger

server:
	go run cmd/server/main.go -port 8080

rest: 
	go run cmd/server/main.go -port 8081 -type rest -endpoint 0.0.0.0:8080