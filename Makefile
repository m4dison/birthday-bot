# Установите переменные для подключения к базе данных
DB_USER=botadmin
DB_PASSWORD=admin
DB_NAME=botdatabase
DB_HOST=localhost
DB_PORT=5432
DB_SSLMODE=disable

# Переменная для строки подключения к базе данных
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Путь к директории с миграциями
MIGRATIONS_DIR=migrations

.PHONY: migrate-up

# Команда для накатывания миграций
migrate-up:
    migrate -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up


