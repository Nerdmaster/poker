.PHONY: example test benchmark

example:
	go build -o bin/poker ./cmd/poker

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -timeout=10m ./...
	go tool cover -html coverage.txt -o coverage.html

benchmark:
	go test -bench=. -benchtime 10s
