FROM golang:1.23 AS builder

RUN mkdir -p /usr/src/simple-auth-service

WORKDIR /usr/src/simple-auth-service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/simple-auth-service ./cmd

FROM alpine:3.20

RUN mkdir -p /usr/local/bin

COPY --from=builder /usr/src/simple-auth-service/build/simple-auth-service /usr/local/bin/simple-auth-service

RUN chmod +x /usr/local/bin/simple-auth-service

EXPOSE 80

CMD ["/usr/local/bin/simple-auth-service"]
