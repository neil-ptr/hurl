build:
	go build

run: build
	go run hurl.go

test: build
	go test ./...
