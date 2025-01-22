package order

import (
	"context"
	"log"
	"time"

	pb "github.com/azizkhan030/go-grpc-graphql/order/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error) {
	protoProducts := []*pb.PostOrderRequest_OrderProduct{}

	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}
	r, err := c.service.PostOrder(
		ctx,
		&pb.PostOrderRequest{
			AccountId: accountId,
			Products:  protoProducts,
		},
	)

	if err != nil {
		return nil, err
	}

	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)

	return &Order{
		ID:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		Products:   products,
		AccountID:  newOrder.AccountId,
	}, nil
}

func (c *Client) GetOrdersFprAccount(ctx context.Context, accountId string) ([]Order, error) {
	r, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrderForAccountRequest{
		AccountId: accountId,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	orders := []Order{}

	for _, orderProto := range r.Orders {
		newOrder := Order{
			ID:         orderProto.Id,
			TotalPrice: orderProto.TotalPrice,
			AccountID:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)
		products := []OrderedProduct{}

		for _, productProto := range orderProto.Products {
			products = append(products, OrderedProduct{
				ID:          productProto.Id,
				Quantity:    productProto.Quantity,
				Name:        productProto.Name,
				Description: productProto.Description,
				Price:       productProto.Price,
			})
		}
		newOrder.Products = products
		orders = append(orders, newOrder)
	}

	return orders, nil
}
