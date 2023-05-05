BINARY_NAME=ports-app

run:
	go run -ldflags="-s -w" cmd/main.go

build:
	go vet cmd/main.go
	go build -o ${BINARY_NAME} -ldflags="-s -w" cmd/main.go

clean:
	rm -f ${BINARY_NAME}

test:
	go test -cover ./... -race


proto:
	protoc --go_opt=paths=source_relative --proto_path=./internal/infra/db/proto --go_out=./internal/infra/db/proto ./internal/infra/db/proto/port.proto