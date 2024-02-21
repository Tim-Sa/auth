package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/Tim-Sa/auth/pkg/auth_v1"
)

const grpcPort = 50051

func randIndex() int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0
	}
	return nBig.Int64()
}

type server struct {
	desc.UnimplementedAuthV1Server
}

func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
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

func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	name := req.Info.GetName()
	email := req.Info.GetEmail()
	password := req.Info.GetPassword()
	passwordConfirm := req.Info.GetPasswordConfirm()
	role := req.Info.GetRole()

	fmt.Println(color.GreenString("New user created:"))
	log.Printf("name - %s ", name)
	log.Printf("email - %s ", email)
	log.Printf("password - %s ", password)
	log.Printf("password confirm - %s ", passwordConfirm)
	log.Printf("role - %v ", role)

	return &desc.CreateResponse{
		Id: randIndex(),
	}, nil
}

func (s *server) Update(_ context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	name := req.GetName()
	email := req.GetEmail()
	role := req.GetRole()

	log.Printf("User with id - %d was updated", id)

	if name != nil {
		log.Printf("name - %s ", name)
	}

	if email != nil {
		log.Printf("email - %s ", email)
	}

	log.Printf("role - %v ", role)

	empty := emptypb.Empty{}
	return &empty, nil
}

func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("User id: %d was removed", req.GetId())

	empty := emptypb.Empty{}
	return &empty, nil
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
