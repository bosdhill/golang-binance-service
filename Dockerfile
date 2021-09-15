FROM golang:1.16 AS build-stage

WORKDIR /go/src/
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o /tmp/server 

FROM scratch

ARG PORT=5000

WORKDIR /
COPY --from=build-stage /tmp/server /server
EXPOSE ${PORT}

CMD ["./server"]