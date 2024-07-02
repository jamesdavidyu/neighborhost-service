build:
	@go build -o bin/neighborhost-service ./app.go

test:
	@go test -v ./...

run: build
	@./bin/neighborhost-service

migration:
	@migrate create -ext sql -dir cmd/model/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/model/migrate/main.go up

migrate-down:
	@go run cmd/model/migrate/main.go down