# Defaults
DSN?='postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable'

# Run migrations
# Usage: make run-migrations DSN=<data_source_name>
# Example: make run-migrations DSN=postgres://user:password@localhost:5432/dbname?sslmode=disable
# Note: DSN is optional, see defaults for more info
.PHONY run-migrations:
run-migrations:
	@echo "Running migrations..."
	@migrate -path "./db/migrations" -database $(DSN) -verbose up

# Create a new migration file
# Usage: make create-migration NAME=<migration_name>
# Example: make create-migration NAME=users_table
.PHONY create-migration:
create-migration:
	@echo "Creating migration..."
	@migrate create -ext sql -dir db/migrations -seq $(NAME)
