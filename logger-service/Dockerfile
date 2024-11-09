FROM golang:1.23 AS builder

RUN mkdir -p /usr/src/simple-logger-service

WORKDIR /usr/src/simple-logger-service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/simple-logger-service ./cmd

FROM alpine:3.20

RUN mkdir -p /usr/local/bin

COPY --from=builder /usr/src/simple-logger-service/build/simple-logger-service /usr/local/bin/simple-logger-service

RUN chmod +x /usr/local/bin/simple-logger-service

EXPOSE 80 5001 50001

CMD ["/usr/local/bin/simple-logger-service"]
