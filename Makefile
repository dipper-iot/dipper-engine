GO111MODULE=on

build:
	go build

test: build
	go test -cover ./... -covermode=count -coverprofile=coverage.out

test-coverage: test
	go tool cover -html=coverage.out