FROM golang:1.16 AS build-stage

WORKDIR /go/src/golang-binance-service
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o /tmp/server \
    && cp .env /tmp/

FROM scratch

ARG PORT=5000

WORKDIR /
COPY --from=build-stage /tmp/server /server
COPY --from=build-stage /tmp/.env /.env
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE ${PORT}

CMD ["./server"]