test:
	go test -coverprofile=coverage.out ./... -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
build:
	go build -o bin/govm cmd/govm/main.go
install: build
	mkdir -p ~/.govm/bin
	cp bin/govm ~/.govm/bin/govm
