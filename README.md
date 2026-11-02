# Metrics

Go API with PostgreSQL. Install Docker with the Compose plugin and GNU Make.

```sh
make migrate-up
make up
```

This creates `.env` from `env_example` if needed, builds `Dockerfile.dev`, and
starts the API with Air after PostgreSQL is ready. The API listens at
http://localhost:8080 and Swagger
is available at http://localhost:8080/swagger/index.html when `Environment=local`.

Follow development logs with:

```sh
make logs SERVICE=app
```

The single `compose.yaml` uses `Dockerfile.dev` and `.air.toml`. Edits to Go files,
SQL migrations, `go.mod`, or `go.sum` rebuild and restart the API automatically.
The source directory is mounted read-only; temporary binaries and Go caches use
Docker volumes. `make up` waits until the API is ready. Stop it with `make down`.
`make dev` and `make dev-down` are aliases for these commands.

`Dockerfile` builds the compiled runtime image without Air or source mounts:

```sh
docker build -f Dockerfile -t metrics:local .
```

Restart development mode after editing `.air.toml`. For local execution without
the API container, start PostgreSQL with `make db-up`, then run
`DB_HOST=127.0.0.1 DB_PORT=5432 go run github.com/air-verse/air@v1.63.0 -c .air.toml`.
Use your configured `POSTGRES_PORT` for `DB_PORT` if different.

Use `make help` to list commands. Common commands:

```sh
make logs SERVICE=app
make ps
make test
make db-shell
make down
```

## Monitoring

`make up` also starts Prometheus and Grafana:

- Prometheus: http://localhost:9090 (scrape status under Targets).
- Grafana: http://localhost:3000, default login `admin` / `admin`.
- API metrics: http://localhost:8080/metrics.

Grafana automatically provisions the Prometheus data source and the **Metrics API**
dashboard in the **Metrics** folder. Panels show availability, request rate, server
errors, p95 latency, goroutines, and memory. Requests and latency use route patterns
to avoid a separate series for each record ID. Prometheus scrapes every 15 seconds;
rate charts need at least two scrapes. Air restarts reset API counters.

Set `PROMETHEUS_PORT`, `GRAFANA_PORT`, `GRAFANA_ADMIN_USER`, and
`GRAFANA_ADMIN_PASSWORD` in `.env` to override defaults. Grafana admin credentials
are applied when its database is first initialized. Prometheus data (15-day
retention) and Grafana settings persist in named volumes across `make down`.

Configuration and dashboard files live in `monitoring/`. Restart Prometheus after
editing its scrape configuration. This setup collects API and Go process metrics;
it does not include a PostgreSQL or Docker host exporter.

Generate 100 API requests for the dashboard (requires Python 3):

```sh
make load-test
# Custom API port or a faster run:
make load-test BASE_URL=http://localhost:8081 DELAY=0.1
```

The script creates one warehouse and product, then cycles through warehouse/product
listing, order creation, validation failures, duplicate conflicts, and a missing
route. Setup counts toward the 100 requests. Expected statuses are 200, 201, 400,
404, and 409; the script reports unexpected responses and exits non-zero. Each run
uses unique test names and leaves its warehouse, product, and orders in the database.
The default delay spreads traffic over about 35 seconds. Allow another scrape
interval before inspecting charts; the server-error panel stays at zero unless
the API returns 5xx responses.

`make down` preserves PostgreSQL data in a named volume. To intentionally delete
the database and start fresh, run `docker compose down --volumes`.

Edit `.env` to change `DB_USER`, `DB_PASS`, or `DB_NAME`. PostgreSQL applies these
values only when initializing an empty volume. `APP_PORT` (default `8080`) and
`POSTGRES_PORT` (default `5432`) control the ports exposed on localhost. Inside
Compose, the API always connects to `postgres:5432` regardless of `DB_HOST` and
`DB_PORT` in `.env`.

To run the API directly with a local Go toolchain matching `go.mod`:

```sh
make migrate-up
DB_HOST=127.0.0.1 DB_PORT=5432 go run ./cmd/api
```

Set `DB_PORT` to your `POSTGRES_PORT` if you changed the published database port.
The API loads the remaining settings from `.env`.

PostgreSQL creates the configured database on first startup. Apply schema changes
explicitly with `make migrate-up` before starting the API and after adding
migrations. The API does not run migrations or GORM AutoMigrate on startup.

Migration files live in `migrations/` at the root of this repository:

- `000001_orders.up.sql` / `.down.sql`: initial orders table; adopts the earlier
  AutoMigrate table if it already exists.
- `000002_warehouses.up.sql` / `.down.sql`: warehouses table and order foreign key.
- `000003_products.up.sql` / `.down.sql`: products belong to warehouses; orders
  reference products instead of warehouses. Existing order names become products
  in their original warehouses, preserving order quantities and links.
- `000004_product_quantity.up.sql` / `.down.sql`: product quantity, defaulting to
  zero for existing products, with a database constraint against negative values.

Each migration runs in a transaction and is recorded in `schema_migrations`.
A PostgreSQL advisory lock prevents concurrent migration runs. The SQL files are
embedded in the migration binary; `make migrate-up` rebuilds it before running.

```sh
make migrate-up     # Start PostgreSQL and apply pending migrations
docker compose stop app
make migrate-down   # Revert only the latest migration
```

Rolling back the first migration drops the orders table and its data. Rolling back the warehouse
migration restores warehouse names on orders and drops the warehouses table.
Standalone warehouse records without orders are lost in that rollback.
Rolling back the product migration restores order names and warehouse IDs, then
drops products. Product records without orders are lost in that rollback.

Create a warehouse first, or list existing warehouses with
`GET /api/v1/warehouses`:

```sh
curl -X POST http://localhost:8080/api/v1/warehouses \
  -H 'Content-Type: application/json' \
  -d '{"name":"Yerevan warehouse"}'
```

Create a product using the warehouse ID returned above, or list existing products
with `GET /api/v1/products`:

```sh
curl -X POST http://localhost:8080/api/v1/products \
  -H 'Content-Type: application/json' \
  -d '{"name":"Office supplies","warehouse_id":1,"quantity":100}'
```

Product responses contain `id`, `name`, `warehouse_id`, and `quantity`.
Product quantity is a non-negative integer and defaults to zero when omitted.
Order creation records the ordered quantity without deducting product quantity.
Product names must be
non-empty and unique within each warehouse. Unknown warehouses return HTTP 404;
duplicate product names in the same warehouse return HTTP 409.

Create an order using the product ID returned above:

```sh
curl -X POST http://localhost:8080/api/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{"product_id":1,"quantity":10}'
```

Successful requests return HTTP 201 with the saved order ID (`{"id":1}`).
`quantity` and `product_id` must be positive integers. Unknown product IDs
return HTTP 404. Orders derive their warehouse through the product relationship.
Warehouse names must be non-empty and unique; duplicates return
HTTP 409. Database failures return HTTP 500 without exposing database details.

Request values are validated once by Gin's `ShouldBindJSON`, using DTO binding
tags. The custom `notblank` rule in `internal/helpers/validation.go` is registered
at startup and rejects whitespace-only names. Services trim accepted names;
repositories handle database constraints such as unknown IDs and duplicates.
Code that calls services directly must validate its request DTO first.

Validation errors return HTTP 400 with messages keyed by JSON field name:

```json
{
  "error": "validation failed",
  "fields": {
    "name": "must not be blank",
    "warehouse_id": "is required"
  }
}
```

Schema validation reports all invalid fields together. JSON parsing or type
errors must be fixed before schema validation runs. Other API errors retain
the `error` message and omit `fields`.

The warehouse migration converts existing order warehouse names into distinct
warehouse records and links the orders by ID. Blank legacy names map to an
`Unassigned` warehouse. Warehouse names are trimmed during conversion. This does
not track warehouse inventory or capacity.

Run `make swagger` with Go installed to regenerate `docs/docs.go`,
`docs/swagger.json`, and `docs/swagger.yaml` from the API annotations. The CLI
version matches the Swagger library pinned in `go.mod`. Rebuild the running API
with `make up` after changing code or generating documentation.
