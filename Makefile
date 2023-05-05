BINARY_NAME=ports-app

run:
	go run -ldflags="-s -w" cmd/main.go

dep:
	go get .

vet:
	go vet ./...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME} -ldflags="-s -w" cmd/main.go

clean:
	rm -f ${BINARY_NAME}
	rm -rf *coverage*
	rm -rf *test-results*

test:
	go test -cover ./... -race

test-coverage:
	go test ./... -covermode=atomic -coverpkg=./... -coverprofile unit-test-coverage.out -race -json > unit-test-results.json
	go tool cover -html unit-test-coverage.out -o unit-test-coverage.html

proto:
	protoc --go_opt=paths=source_relative --proto_path=./internal/infra/db/proto --go_out=./internal/infra/db/proto ./internal/infra/db/proto/port.proto

docker-build: build
	docker build -t fabricioandreis/ports-app:latest .

docker-run: 
	docker run --detach --env-file ./local.env --volume ${PWD}/ports.json:/data/ports.json --network=host fabricioandreis/ports-app:latest

docker-brun: docker-build docker-run

docker-push: docker-build
	docker push fabricioandreis/ports-app

up: docker-build
	docker compose up --abort-on-container-exit --exit-code-from app --remove-orphans

acceptance-tests:
	docker compose up --abort-on-container-exit --exit-code-from app --remove-orphans