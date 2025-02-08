test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
build:
	go build -o bin/govm cmd/govm/main.go
install: build
	mkdir -p ~/.govm/bin
	cp bin/govm ~/.govm/bin/govm
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/govm cmd/govm/main.go
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/govm cmd/govm/main.go
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/govm cmd/govm/main.go
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin/arm64/govm cmd/govm/main.go
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64
