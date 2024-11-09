FROM golang:1.23 AS builder

RUN mkdir -p /usr/src/simple-gateway-service

WORKDIR /usr/src/simple-gateway-service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/simple-gateway-service ./cmd

FROM alpine:3.20

RUN mkdir -p /usr/local/bin

COPY --from=builder /usr/src/simple-gateway-service/build/simple-gateway-service /usr/local/bin/simple-gateway-service

RUN chmod +x /usr/local/bin/simple-gateway-service

EXPOSE 80

CMD ["/usr/local/bin/simple-gateway-service"]
