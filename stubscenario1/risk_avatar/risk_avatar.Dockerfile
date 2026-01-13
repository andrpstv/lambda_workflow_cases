FROM golang:1.24-alpine3.22 AS builder
WORKDIR /app
COPY go.mod .
COPY . .
ENV CGO_ENABLED=0
RUN go build -o build/main ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder ./app/build .
RUN chmod +x main
CMD ["./main"]