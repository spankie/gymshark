package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	// postgres driver
	_ "github.com/lib/pq"
	"github.com/spankie/gymshark/database/models"
)

type Service interface {
	Health(ctx context.Context) (string, error)
	CreateOrder(ctx context.Context, order *models.Order, orderShipping []*models.OrderShipping) error
	GetOrder(ctx context.Context, id int) (*models.Order, error)
	GetAvailableShippingPacks(ctx context.Context) ([]models.ShippingPack, error)
	GetOrderShipping(ctx context.Context) ([]models.OrderShipping, error)
}

type postgresService struct {
	db *sql.DB
}

// getDBConnection returns a connection to the postgres database
func getDBConnection(port int, host, username, password, dbname string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return db, nil
}

// NewPostgresDBService creates a new postgres database connection and returns an
// implementation of db service
func NewPostgresDBService(port int, host, username, password, dbname string) (Service, error) {
	db, err := getDBConnection(port, host, username, password, dbname)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	err = MigrateDb(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname))
	if err != nil {
		return nil, fmt.Errorf("could not migrate database: %w", err)
	}

	return &postgresService{
		db: db,
	}, nil
}

// Health checks if the database is up
func (ps *postgresService) Health(ctx context.Context) (string, error) {
	err := ps.db.PingContext(ctx)
	if err != nil {
		return "", fmt.Errorf("postgres db down: %w", err)
	}

	return "postgres is healthy", nil
}

// CreateOrder inserts an order into the database
func (ps *postgresService) CreateOrder(ctx context.Context, order *models.Order, orderShipping []*models.OrderShipping) error {
	tx, err := ps.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("unable to start db transaction: %w", err)
	}

	query := `INSERT INTO orders (id, number_of_items) VALUES (DEFAULT, $1) RETURNING id`
	row := tx.QueryRowContext(ctx, query, order.NumberOfItems)
	err = row.Scan(&order.ID)
	if err != nil {
		return errors.Join(fmt.Errorf("could not insert order: %w", err), rollback(tx))
	}

	queryOrderShipping := `INSERT INTO order_shipping
	(id, order_id, pack_size, shipping_pack_quantity)
	VALUES (DEFAULT, $1, $2, $3) RETURNING id`
	for k, v := range orderShipping {
		row := tx.QueryRowContext(ctx, queryOrderShipping, order.ID, v.PackSize, v.ShippingPackQuantity)
		err := row.Scan(&orderShipping[k].ID)
		if err != nil {
			return errors.Join(fmt.Errorf("could not insert order_shipping: %w", err), rollback(tx))
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit db transaction: %w", err)
	}

	return nil
}

func rollback(tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil {
		slog.Error("error rolling back transaction", "error", err)
		return err
	}

	return nil
}

func (ps *postgresService) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	query := `SELECT id, number_of_items, created_at, updated_at FROM orders where id = $1`
	row := ps.db.QueryRowContext(ctx, query, id)

	var order models.Order
	err := row.Scan(&order.ID, &order.NumberOfItems, &order.CreatedAt, &order.UpdateAt)
	if err != nil {
		return nil, fmt.Errorf("could not get order: %w", err)
	}

	// fetch the order shipping
	shippingQuery := `SELECT id, pack_size, shipping_pack_quantity FROM order_shipping WHERE order_id = $1 ORDER BY pack_size DESC`
	rows, err := ps.db.QueryContext(ctx, shippingQuery, order.ID)
	if err != nil {
		return nil, fmt.Errorf("error finding shipping details for order: %v", err)
	}

	for rows.Next() {
		shipping := models.OrderShipping{}
		err := rows.Scan(&shipping.ID, &shipping.PackSize, &shipping.ShippingPackQuantity)
		if err != nil {
			return nil, fmt.Errorf("could not get order shipping information: %v", err)
		}
		order.Shipping = append(order.Shipping, shipping)
	}

	return &order, nil
}

func (ps *postgresService) GetAvailableShippingPacks(ctx context.Context) ([]models.ShippingPack, error) {
	query := `SELECT id, quantity, created_at, updated_at FROM shipping_packs ORDER BY quantity DESC`
	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error query db for shipping packs: %w", err)
	}

	var packs []models.ShippingPack
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error scanning columns from shipping pack: %w", err)
		}
		var pack models.ShippingPack
		err := rows.Scan(&pack.ID, &pack.Quantity, &pack.CreatedAt, &pack.UpdateAt)
		if err != nil {
			return nil, fmt.Errorf("could not get order: %w", err)
		}
		packs = append(packs, pack)
	}

	return packs, nil
}

func (ps *postgresService) GetOrderShipping(ctx context.Context) ([]models.OrderShipping, error) {
	query := `SELECT id, order_id, pack_size, shipping_pack_quantity, created_at, updated_at FROM order_shipping ORDER BY created_at DESC`
	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting order shipping from db: %v", err)
	}
	orderShipping := []models.OrderShipping{}
	for rows.Next() {
		s := models.OrderShipping{}
		err := rows.Scan(&s.ID, &s.OrderID, &s.PackSize, &s.ShippingPackQuantity, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orderShipping = append(orderShipping, s)
	}

	return orderShipping, nil
}
