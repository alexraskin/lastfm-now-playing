FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .

COPY handlers/ ./handlers/

COPY models/ ./models/

COPY service/ ./service/

COPY utils/ ./utils/

COPY templates/ ./templates/

RUN CGO_ENABLED=0 GOOS=linux go build -o lastfm-api main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/lastfm-api .
COPY --from=builder /app/templates ./templates

EXPOSE 3000

CMD ["./lastfm-api"] 