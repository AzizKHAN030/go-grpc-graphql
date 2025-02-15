package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/azizkhan030/go-grpc-graphql/catalog/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	serv := grpc.NewServer()

	pb.RegisterCatalogServiceServer(serv, &grpcServer{
		pb.UnimplementedCatalogServiceServer{},
		s,
	})

	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.PostProductResponse{
		Product: productOut(p),
	}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: productOut(p),
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res []Product
	var err error

	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductsByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.GetProductsResponse{
		Products: mapProductsToProductsResponse(res),
	}, nil
}

func mapProductsToProductsResponse(p []Product) []*pb.Product {
	products := []*pb.Product{}

	for _, product := range p {
		products = append(products, productOut(&product))
	}

	return products
}

func productOut(p *Product) *pb.Product {
	return &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
	}
}
