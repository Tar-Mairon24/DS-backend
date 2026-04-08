include .env
export

MIGRATE=~/go/bin/migrate
DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&loc=Local

migrate-up:
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" down 1

migrate-version: 
	$(MIGRATE) -path ./migrations -database "$(DB_URL)" version

migrate-create:
	$(MIGRATE) create -ext sql -dir ./migrations -seq $(name)

create-backend:
	@docker network create ds-network || true
	@docker compose up -d --build
	@sleep 5
	@$(MAKE) migrate-up
	@bash ./scripts/seed_auth.sh
	@bash ./scripts/seed_data.sh
	@bash ./scripts/seed_images.sh
	@echo "Backend created and seeded!"

seed:
	@echo "Starting database seeding..."
	@echo "1. Authenticating..."
	@bash ./scripts/seed_auth.sh
	@echo ""
	@echo "2. Seeding users and properties..."
	@bash ./scripts/seed_data.sh
	@echo ""
	@echo "3. Seeding images..."
	@bash ./scripts/seed_images.sh
	@echo ""
	@echo "Database seeding complete!"

restart-database: 
	@docker compose down -v
	@docker compose up -d --build
	@sleep 5
	@$(MAKE) migrate-up
	@bash ./scripts/seed_auth.sh
	@bash ./scripts/seed_data.sh
	@bash ./scripts/seed_images.sh
	@echo "Database reset and seeded!"