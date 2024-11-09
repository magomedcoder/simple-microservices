FROM golang:1.23 AS builder

RUN mkdir -p /usr/src/simple-mailer-service

WORKDIR /usr/src/simple-mailer-service

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/simple-mailer-service ./cmd

FROM alpine:3.20

RUN mkdir -p /usr/local/bin /usr/local/share/simple-mailer-service/templates

COPY --from=builder /usr/src/simple-mailer-service/build/simple-mailer-service /usr/local/bin/simple-mailer-service

COPY --from=builder /usr/src/simple-mailer-service/templates /usr/local/share/simple-mailer-service/templates

RUN chmod +x /usr/local/bin/simple-mailer-service

EXPOSE 80

CMD ["/usr/local/bin/simple-mailer-service"]
