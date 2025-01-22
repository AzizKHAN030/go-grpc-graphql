package account

import (
	"context"

	accountpb "github.com/azizkhan030/go-grpc-graphql/account/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	connection *grpc.ClientConn
	service    accountpb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	connection, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	c := accountpb.NewAccountServiceClient(connection)

	return &Client{connection, c}, nil
}

func (c *Client) Close() {
	c.connection.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(ctx, &accountpb.PostAccountRequest{Name: name})

	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(ctx, &accountpb.GetAccountRequest{Id: id})

	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(ctx, &accountpb.GetAccountsRequest{Skip: skip, Take: take})

	if err != nil {
		return nil, err
	}

	var accounts []Account

	for _, acc := range r.Accounts {
		accounts = append(accounts, Account{
			ID:   acc.Id,
			Name: acc.Name,
		})
	}

	return accounts, nil
}
