.PHONY: test coverage lint run clean

test:
	golangci-lint run
	go fmt ./...
	go test -v ./...
	make coverage

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

run:
	make clean
	go mod tidy
	make test
	go run ./cmd/main.go

clean:
	rm -f coverage.out coverage.html
