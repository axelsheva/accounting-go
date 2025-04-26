.PHONY: up down run generate new-entity tidy

# Start PostgreSQL in Docker
up:
	docker-compose up -d
	sleep 1

# Stop containers
down:
	docker-compose down --volumes

# Run application
run:
	go run main.go

# Generate Ent code
generate:
	go generate ./ent

# Create new entity (usage: make new-entity NAME=EntityName)
new-entity:
	go run -mod=mod entgo.io/ent/cmd/ent new $(NAME)

# Update dependencies
tidy:
	go mod tidy

# Start PostgreSQL and application
start: up run
