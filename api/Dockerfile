FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN CGO_ENABLED=0 go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app ./main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]