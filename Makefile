
CONTAINER_NAME=postgres
DB_NAME=simple_bank
MIGRATION_NAME=init_schema
MIGRATE_VERSION=

postgres:
	docker run --name $(CONTAINER_NAME) -e POSTGRES_PASSWORD=1337 -e POSTGRES_USER=root -p 5433:5432 -d postgres:12-alpine

createdb:
	docker exec -it $(CONTAINER_NAME) createdb --username=root --owner=root $(DB_NAME)

dropdb:
	docker exec -it $(CONTAINER_NAME) dropdb $(DB_NAME)

.PHONY: createdb

connect:
	docker exec -it $(CONTAINER_NAME)  /bin/bash

newmigrate:
	 migrate create -ext sql -dir db/migration -seq $(MIGRATION_NAME)

migrateup:
	 migrate -path /Users/grishu/simplebank/db/migration -database "postgresql://root:1337@127.0.0.1:5433/simple_bank?sslmode=disable" -verbose up $(MIGRATE_VERSION)

migratedown:
	 migrate -path /Users/grishu/simplebank/db/migration -database "postgresql://root:1337@127.0.0.1:5433/simple_bank?sslmode=disable" -verbose down $(MIGRATE_VERSION)

sqlc:
	sqlc generate

test:
	go test -cover ./...

