package catalog

import (
	"context"
	"strings"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	GetProductsByIDs(ctx context.Context, ids string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{r}
}

func (s *catalogService) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	p := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().String(),
	}

	if err := s.repository.PutProduct(ctx, *p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	p, err := s.repository.GetProductByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *catalogService) GetProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.ListProducts(ctx, skip, take)
}

func (s *catalogService) GetProductsByIDs(ctx context.Context, ids string) ([]Product, error) {
	idsSlice := strings.Split(ids, ",")
	p, err := s.repository.ListProductsWithIDs(ctx, idsSlice)

	if err != nil {
		return nil, err
	}

	return p, nil

}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.SearchProducts(ctx, query, skip, take)
}
