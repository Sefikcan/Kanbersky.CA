.PHONY: migrate migrate_down migrate_up migrate_version docker prod local swaggo

# ==================================================================
# Migration

force:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations down 1

# ===================================================================
# Tools Commands

swaggo:
	echo "Starting swagger generating"
	swag init -g ./cmd/server/main.go -o ./docs

# ===================================================================
# Main

run:
	go run ./cmd/server/main.go

build:
	go build ./cmd/server/main.go