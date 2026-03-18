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