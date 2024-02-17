package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"crypto/rand"
	"net"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"

	desc "github.com/Tim-Sa/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

func RandIndex() int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0
	}
	return nBig.Int64()
}

type server struct {
	desc.UnimplementedAuthV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	var role desc.Role
	if gofakeit.Bool() {
		role = desc.Role_ADMIN
	} else {
		role = desc.Role_USER
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id:        req.GetId(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Role:      role,
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	name := req.Info.GetName()
	email := req.Info.GetEmail()
	password := req.Info.GetPassword()
	password_confirm := req.Info.GetPasswordConfirm()
	role := req.Info.GetRole()

	fmt.Println(color.GreenString("New user created:"))
	log.Printf("name - %s ", name)
	log.Printf("email - %s ", email)
	log.Printf("password - %s ", password)
	log.Printf("password_confirm - %s ", password_confirm)
	log.Printf("role - %v ", role)

	return &desc.CreateResponse{
		Id: RandIndex(),
	}, nil
}

func main() {
	fmt.Println(color.GreenString("Auth service start"))

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", listen.Addr())

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}