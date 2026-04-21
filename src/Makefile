.PHONY: run-farm lint test docker-up docker-down

run-farm:
	go run ./apps/farm-service/cmd/

lint:
	golangci-lint run ./...

test:
	go test ./...

docker-up:
	docker compose -f ../deployments/docker-compose.yaml up -d postgres redis

docker-down:
	docker compose -f ../deployments/docker-compose.yaml down
