FROM golang:1.20-alpine

WORKDIR /app

COPY --chmod=0755 ./internal/tests/entrypoint.sh ./

COPY . ./

ENTRYPOINT [ "/app/entrypoint.sh" ]