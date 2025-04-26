.PHONY: up down run generate new-entity tidy

# Запуск PostgreSQL в Docker
up:
	docker-compose up -d

# Остановка контейнеров
down:
	docker-compose down

# Запуск приложения
run:
	go run main.go

# Генерация кода Ent
generate:
	go generate ./ent

# Создание новой сущности (использование: make new-entity NAME=ИмяСущности)
new-entity:
	go run -mod=mod entgo.io/ent/cmd/ent new $(NAME)

# Обновление зависимостей
tidy:
	go mod tidy

# Запуск PostgreSQL и приложения
start: down up run
