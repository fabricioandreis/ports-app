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
