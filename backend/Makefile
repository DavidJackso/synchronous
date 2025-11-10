run:
	go run cmd/app/main.go

# Миграции с goose
# Использование: DB_DSN="postgres://user:password@localhost:5432/database?sslmode=disable" make migrate-up
DB_DSN ?= postgres://user:password@localhost:5432/database?sslmode=disable

migrate-up:
	goose -dir migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir migrations postgres "$(DB_DSN)" down

migrate-status:
	goose -dir migrations postgres "$(DB_DSN)" status

migrate-create:
	@read -p "Введите название миграции: " name; \
	goose -dir migrations create $$name sql

migrate-reset:
	goose -dir migrations postgres "$(DB_DSN)" reset