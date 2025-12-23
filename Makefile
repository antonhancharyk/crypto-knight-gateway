.PHONY: build docker run test
build:
	go build -o bin/gateway ./cmd/gateway
docker:
	docker build -t gateway -f Dockerfile .
run:
	./bin/gateway
test:
	go test ./... -v