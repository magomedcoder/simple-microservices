FROM golang:1.23 AS builder

RUN mkdir -p /usr/src/simple-listener-service

WORKDIR /usr/src/simple-listener-service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/simple-listener-service ./cmd

FROM alpine:3.20

RUN mkdir -p /usr/local/bin

COPY --from=builder /usr/src/simple-listener-service/build/simple-listener-service /usr/local/bin/simple-listener-service

RUN chmod +x /usr/local/bin/simple-listener-service

EXPOSE 80

CMD ["/usr/local/bin/simple-listener-service"]
