test:
	go test -v ./...
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/govm cmd/govm/main.go
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/govm cmd/govm/main.go
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/govm cmd/govm/main.go
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin/arm64/govm cmd/govm/main.go
build-all:
	make build-linux-amd64
	make build-linux-arm64
	make build-darwin-amd64
	make build-darwin-arm64
