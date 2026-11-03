# Metrics

A Go API for managing warehouses, products, and orders with PostgreSQL.
Includes Prometheus metrics, Loki, and Grafana monitoring. Requires Docker
with the Compose plugin and GNU Make.

```sh
make migrate-up    # Start PostgreSQL and apply migrations
make up            # Start the API and monitoring services
make down          # Stop services and preserve data
make help          # List available commands
make env           # Create .env from env_example if missing
make build         # Build the application image
make restart       # Restart the API
make logs          # Follow all service logs
make logs SERVICE=app # Follow API logs
make ps            # Show service status
make config        # Validate Compose configuration
make db-up         # Start PostgreSQL only
make db-shell      # Open the PostgreSQL shell
make test          # Run Go tests in Docker
make load-test     # Send 100 API requests (requires Python 3)
make swagger       # Generate Swagger documentation (requires Go)
make migrate-down # Revert the latest migration (stop the API first)
make dev           # Alias for make up
make dev-down      # Alias for make down
```
