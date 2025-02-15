package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type dbRepository struct {
	db *sql.DB
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

func NewDbRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &dbRepository{db}, nil
}

func (r *dbRepository) Close() {
	r.db.Close()
}

func (r *dbRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, p := range o.Products {
		if _, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity); err != nil {
			return err
		}
	}

	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *dbRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
		o.id,
		o.created_atm
		o.account_id,
		o.total_price::money::numeric::float8,
		op.product_id,
		op.quantity
		FROM orders o JOIN order_products op ON(o.id = op.order_id)
		ORDER BY o.id`,
		accountID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	products := []OrderedProduct{}
	order := Order{}

	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrder.ID != "" && lastOrder.ID != order.ID {
			newOrder := Order{
				ID:         lastOrder.ID,
				AccountID:  lastOrder.AccountID,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   lastOrder.Products,
			}

			orders = append(orders, newOrder)
			products = []OrderedProduct{}
		}
		products = append(products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})

		*lastOrder = *&order
	}

	if lastOrder != nil {
		newOrder := Order{
			ID:         lastOrder.ID,
			AccountID:  lastOrder.AccountID,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   lastOrder.Products,
		}

		orders = append(orders, newOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
