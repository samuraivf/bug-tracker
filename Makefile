run:
	go run ./cmd/bug-tracker/main.go

build:
	go build -o bin/main ./cmd/bug-tracker/main.go

test:
	go test -v -count=1 ./...

test100:
	go test -v -count=100 ./...

race:
	go test -v -race -count=1 ./...

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
