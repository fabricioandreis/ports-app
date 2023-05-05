FROM gcr.io/distroless/static

WORKDIR /app

COPY ports-app ./

CMD [ "/app/ports-app" ]