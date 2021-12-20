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

build-linux: frontend
	GOOS=linux GOARCH=amd64 go build -o bin/milvus-ops-linux ./cmd/milvus-ops.go

docker: build-linux
	mv ./bin/milvus-ops-linux ./docker/server
	mv ./web/build ./docker/
	cd docker && docker build -t haorenfsa/milvus-ops .