FROM gcr.io/distroless/static

WORKDIR /app

COPY ports-app ./
COPY ports.json /data/

CMD [ "/app/ports-app" ]