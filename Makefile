all: generate test build run

.PHONY: all generate test-backend backend frontend build run

generate: go generate ./...

test-backend: generate
	go test ./...

backend: 
	mkdir -p bin
	go build -o bin/ ./cmd/milvus-ops.go

frontend:
	cd web && yarn install && yarn build

frontend-dev: backend
	cd web && yarn install && yarn start

build: backend frontend

run:
	./bin/milvus-ops