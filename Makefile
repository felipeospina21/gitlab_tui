start:
	@echo "Starting app..."
	go run ./cmd/glabt/glabt.go

build:
	go build ./cmd/glabt/glabt.go

test:
	go test ./...

cover: 
	go test ./... -coverprofile=c.out
	go tool cover -html="c.out"
