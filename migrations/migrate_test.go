package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Run with METRICS_TEST_DSN pointing at a disposable PostgreSQL database.
// Each case uses an isolated schema and removes it when finished.
func TestMigrations(t *testing.T) {
	dsn := os.Getenv("METRICS_TEST_DSN")
	if dsn == "" {
		t.Skip("set METRICS_TEST_DSN to run PostgreSQL migration tests")
	}
	for _, legacy := range []bool{false, true} {
		t.Run(fmt.Sprintf("legacy=%v", legacy), func(t *testing.T) {
			db, err := sql.Open("pgx", dsn)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			db.SetMaxOpenConns(1)
			schema := fmt.Sprintf("migration_test_%d", time.Now().UnixNano())
			exec := func(query string) {
				t.Helper()
				if _, err := db.Exec(query); err != nil {
					t.Fatal(err)
				}
			}
			exec("CREATE SCHEMA " + schema)
			defer db.Exec("DROP SCHEMA " + schema + " CASCADE")
			exec("SET search_path TO " + schema)
			if legacy {
				exec(`CREATE TABLE orders (id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL,
                    created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
                    quantity BIGINT NOT NULL DEFAULT 1, warehouse TEXT NOT NULL DEFAULT '');
                    INSERT INTO orders(name, quantity, warehouse) VALUES
                    ('Old order', 3, ' Yerevan '), ('Second order', 2, 'Yerevan'), ('Unknown', 1, '')`)
			}
			run := func(direction string) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				if err := Run(ctx, db, direction); err != nil {
					t.Fatal(err)
				}
			}
			run("up")
			run("up") // Already-applied migrations must be skipped.
			if legacy {
				var count int
				if err := db.QueryRow(`SELECT count(*) FROM orders o
                        JOIN products p ON p.id=o.product_id JOIN warehouses w ON w.id=p.warehouse_id
                        WHERE (p.name='Old order' AND o.quantity=3 AND w.name='Yerevan')
                        OR (p.name='Second order' AND o.quantity=2 AND w.name='Yerevan')
                        OR (p.name='Unknown' AND o.quantity=1 AND w.name='Unassigned')`).Scan(&count); err != nil || count != 3 {
					t.Fatalf("legacy order links: count=%d, error=%v", count, err)
				}
			}
			if _, err := db.Exec(`INSERT INTO orders(quantity, product_id) VALUES (1, 999999)`); err == nil {
				t.Fatal("foreign key allowed an unknown product")
			}
			if _, err := db.Exec(`INSERT INTO products(name, warehouse_id) VALUES ('Invalid', 999999)`); err == nil {
				t.Fatal("foreign key allowed an unknown warehouse")
			}
			var warehouseID uint
			if err := db.QueryRow(`INSERT INTO warehouses(name) VALUES ('Quantity test') RETURNING id`).Scan(&warehouseID); err != nil {
				t.Fatal(err)
			}
			if _, err := db.Exec(`INSERT INTO products(name, warehouse_id, quantity) VALUES ('Negative quantity', $1, -1)`, warehouseID); err == nil {
				t.Fatal("database allowed negative product quantity")
			}
			var quantity int64
			if err := db.QueryRow(`INSERT INTO products(name, warehouse_id) VALUES ('Default quantity', $1) RETURNING quantity`, warehouseID).Scan(&quantity); err != nil || quantity != 0 {
				t.Fatalf("default quantity=%d, error=%v", quantity, err)
			}
			if _, err := db.Exec(`INSERT INTO products(name, warehouse_id, quantity) VALUES ('Stocked product', $1, 100)`, warehouseID); err != nil {
				t.Fatal(err)
			}
			run("down") // Remove product quantities, preserving products and orders.
			run("down") // Restore order names and warehouse IDs.
			if legacy {
				var count int
				if err := db.QueryRow(`SELECT count(*) FROM orders o JOIN warehouses w ON w.id=o.warehouse_id
                        WHERE o.name='Old order' AND o.quantity=3 AND w.name='Yerevan'`).Scan(&count); err != nil || count != 1 {
					t.Fatalf("product rollback order: count=%d, error=%v", count, err)
				}
			}
			run("down")
			if legacy {
				var name string
				if err := db.QueryRow("SELECT warehouse FROM orders WHERE name='Old order'").Scan(&name); err != nil || name != "Yerevan" {
					t.Fatalf("rollback warehouse name=%q, error=%v", name, err)
				}
			}
			run("up") // Reapply warehouse and product conversion after rollback.
			run("down")
			run("down")
			run("down")
			run("down")
			run("down") // Empty history is a no-op.
			run("up")   // Recreate the whole schema.
		})
	}
}
