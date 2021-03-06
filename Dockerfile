FROM golang:1.18.3-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN \
  unset GOPATH && \
  go build -o main ./main.go


FROM alpine:3.16.0
RUN adduser -h /app -D appuser
WORKDIR /app
COPY --from=builder /app/main .

USER appuser
ENTRYPOINT ["./main"]
EXPOSE 8080/tcp

