FROM gcr.io/distroless/static

WORKDIR /app

COPY --chmod=0755 ports-app ./

ENTRYPOINT [ "executable" ] [ "/app/ports-app" ]