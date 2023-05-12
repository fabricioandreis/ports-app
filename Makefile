BINARY_NAME=ports-app

run:
	go run -ldflags="-s -w" cmd/main.go

dep:
	go get .

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run --enable-all

install-lint:
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin 
	golangci-lint --version

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME} -ldflags="-s -w" cmd/main.go

clean: clear

clear:
	rm -f ${BINARY_NAME}
	rm -rf *coverage*
	rm -rf *test-results*
	go clean -testcache

test:
	go test -cover ./... -skip "TestAcceptance" -race

test-coverage:
	go test ./... -skip "TestAcceptance" -covermode=atomic -coverpkg=./... -coverprofile unit-tests-coverage.out -race -json > unit-test-results.json
	go tool cover -html unit-tests-coverage.out -o unit-tests-coverage.html

proto:
	protoc --go_opt=paths=source_relative --proto_path=./internal/infra/db/proto --go_out=./internal/infra/db/proto ./internal/infra/db/proto/port.proto

docker-build: build
	docker build -t fabricioandreis/ports-app:latest .

docker-run: 
	docker run --detach --env-file ./local.env --volume ${PWD}/ports.json:/data/ports.json --network=host fabricioandreis/ports-app:latest

docker-brun: docker-build docker-run

docker-push: docker-build
	docker push fabricioandreis/ports-app

local: build
	docker compose up --build --exit-code-from app --remove-orphans

local-acceptance-tests: build pipeline-acceptance-tests

# This rule does not build the binary, because during the pipeline it is downloaded from the artifact repository
pipeline-acceptance-tests:
	docker compose --file docker-compose-accept-tests.yaml build
	docker compose --file docker-compose-accept-tests.yaml run acceptance-tests --exit-code-from acceptance-tests --remove-orphans > acceptance-tests-results.log