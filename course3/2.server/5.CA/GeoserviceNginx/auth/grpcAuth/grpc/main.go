package main

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	pb "go-kata/2.server/5.CA/GeoserviceNginx/auth/grpcAuth/grpc/example"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedTokenValidationServiceServer
}

func (s *server) ValidateToken(ctx context.Context, request *pb.TokenRequest) (*pb.TokenResponse, error) {

	tokenString := request.Token

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Замените "secret" на ваш секретный ключ для подписи токена
	})

	if err != nil {
		return &pb.TokenResponse{Valid: false, ErrorMessage: err.Error()}, nil
	}

	if token.Valid {
		return &pb.TokenResponse{Valid: true}, nil
	} else {
		return &pb.TokenResponse{Valid: false, ErrorMessage: "Invalid token"}, nil
	}
}

func main() {
	//запуск gRPC сервера
	lis, err := net.Listen("tcp", ":50051")
	log.Println("gRPC server started on port :50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTokenValidationServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

//	geoService := new(GeoService)
//	rpc.Register(geoService)
//
//	//запуск RPC сервера
//	listener, err := net.Listen("tcp", ":8070")
//	if err != nil {
//		log.Fatal("Error starting RPC server:", err)
//	}
//	log.Println("RPC server started on port 8070")
//
//	//запуск json-RPC сервера
//	listenerJsonRpc, err := net.Listen("tcp", ":8060")
//	if err != nil {
//		log.Fatal("Listen error:", err)
//	}
//	log.Println("JsonRPC server started on port 8060")
//
//	go func() {
//		for {
//			conn, err := listener.Accept()
//			if err != nil {
//				log.Println("Accept error:", err)
//				continue
//			}
//			go func(conn net.Conn) {
//				defer conn.Close()
//				rpc.ServeConn(conn)
//			}(conn)
//		}
//	}()
//
//	go func() {
//		for {
//			connJsonRpc, err := listenerJsonRpc.Accept()
//			if err != nil {
//				log.Println("Accept error:", err)
//				continue
//			}
//			go func(connJsonRpc net.Conn) {
//				defer connJsonRpc.Close()
//				jsonrpc.ServeConn(connJsonRpc)
//			}(connJsonRpc)
//		}
//	}()
//
//	select {} // для предотвращения завершения программы
//}
