package account

import (
	"context"
	"fmt"
	"net"

	accountpb "github.com/azizkhan030/go-grpc-graphql/account/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	accountpb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	server := grpc.NewServer()
	accountpb.RegisterAccountServiceServer(server, &grpcServer{accountpb.UnimplementedAccountServiceServer{}, s})
	reflection.Register(server)

	return server.Serve(lis)
}

func (s *grpcServer) PostAccount(ctx context.Context, r *accountpb.PostAccountRequest) (*accountpb.PostAccountResponse, error) {
	acc, err := s.service.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}

	return &accountpb.PostAccountResponse{
		Account: &accountpb.Account{
			Id:   acc.ID,
			Name: acc.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *accountpb.GetAccountRequest) (*accountpb.GetAccountResponse, error) {
	acc, err := s.service.GetAccount(ctx, r.Id)

	if err != nil {
		return nil, err
	}

	return &accountpb.GetAccountResponse{
		Account: &accountpb.Account{
			Id:   acc.ID,
			Name: acc.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *accountpb.GetAccountsRequest) (*accountpb.GetAccountsResponse, error) {
	accs, err := s.service.GetAccounts(ctx, r.Skip, r.Take)

	if err != nil {
		return nil, err
	}

	accounts := []*accountpb.Account{}

	for _, acc := range accs {
		accounts = append(accounts, &accountpb.Account{
			Id:   acc.ID,
			Name: acc.Name,
		})
	}

	return &accountpb.GetAccountsResponse{Accounts: accounts}, nil
}
