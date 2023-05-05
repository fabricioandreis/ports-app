FROM gcr.io/distroless/static

WORKDIR /app

COPY ports-app ./
COPY ports.json ./

CMD [ "/app/ports-app" ]