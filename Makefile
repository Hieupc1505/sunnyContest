
.PHONY: run env postgres createdb dropdb migrateup migratedown sqlc

env:
	direnv allow

run:
	go run ./cmd/server/main.go || true

postgres:
	docker run --name ${DOCKER_DB_NAME} -p 5432:5432 -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:16-alpine

# Định nghĩa target để tạo database
createdb:
	docker exec -it ${DOCKER_DB_NAME} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

# Định nghĩa target để xóa database
dropdb:
	docker exec -it ${DOCKER_DB_NAME} dropdb ${DB_NAME}

#migrate create -ext sql -dir migrations -seq init_schema
#Định nghĩa target để thực hiện tạo table db/migration
migrateup:
	migrate -path migrations -database "${DB_SOURCE}" -verbose up

# Định nghĩa target để thực hiện xóa db/migration
migratedown:
	migrate -path migrations -database "${DB_SOURCE}" -verbose down

#Định nghĩa sqlc
sqlc:
	sqlc generate

#build
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/app_prod cmd/server/main.go
	@echo "compiled you application with all its assets to a single binary => bin/app_prod"