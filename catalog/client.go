package catalog

import (
	"context"

	pb "github.com/azizkhan030/go-grpc-graphql/catalog/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	connection *grpc.ClientConn
	service    pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	connection, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	c := pb.NewCatalogServiceClient(connection)

	return &Client{
		connection,
		c,
	}, err
}

func (c *Client) Close() {
	c.connection.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	r, err := c.service.PostProduct(ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:    r.Product.Id,
		Name:  r.Product.Description,
		Price: r.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:    r.Product.Id,
		Name:  r.Product.Description,
		Price: r.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, query, ids string, take uint64, skip uint64) ([]*Product, error) {
	r, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Ids:   ids,
		Take:  take,
		Skip:  skip,
		Query: query,
	})

	if err != nil {
		return nil, err
	}

	products := []*Product{}

	for _, p := range r.Products {
		products = append(products, &Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return products, nil
}
