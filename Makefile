BINARY_NAME=ports-app

run:
	go run -ldflags="-s -w" cmd/main.go

build:
	go vet cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME} -ldflags="-s -w" cmd/main.go

clean:
	rm -f ${BINARY_NAME}

test:
	go test -cover ./... -race


proto:
	protoc --go_opt=paths=source_relative --proto_path=./internal/infra/db/proto --go_out=./internal/infra/db/proto ./internal/infra/db/proto/port.proto

docker-build: build
	docker build -t fabricioandreis/ports-app:latest .

docker-run: 
	docker run --detach --env-file ./local.env --network=host fabricioandreis/ports-app

docker-brun: docker-build
	docker run --detach --env-file ./local.env --network=host fabricioandreis/ports-app

docker-push: docker-build
	docker push fabricioandreis/ports-app

up: docker-build
	docker compose up