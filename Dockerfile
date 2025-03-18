FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o lastfm-api main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/lastfm-api .

EXPOSE 3000

CMD ["./lastfm-api"] 