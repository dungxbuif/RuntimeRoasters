.PHONY: proto run-demo lint test docker-up docker-down

BUF=$(shell go env GOPATH)/bin/buf

proto:
	cd ../api && $(BUF) dep update && $(BUF) generate

proto-lint:
	cd ../api && $(BUF) lint

proto-breaking:
	cd ../api && $(BUF) breaking --against '../.git#branch=main,subdir=api'

run-demo:
	go run ./apps/demo-service/cmd/

run-client:
	cd apps/client-app && npm run dev

lint:
	golangci-lint run ./...

test:
	go test ./...

docker-up:
	docker compose -f ../deployments/docker-compose.yaml up -d --remove-orphans

docker-down:
	docker compose -f ../deployments/docker-compose.yaml down
