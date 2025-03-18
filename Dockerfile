FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .

COPY handlers/ ./handlers/

COPY models/ ./models/

COPY service/ ./service/

COPY utils/ ./utils/

RUN CGO_ENABLED=0 GOOS=linux go build -o lastfm-api main.go

FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/lastfm-api .

EXPOSE 3000

CMD ["./lastfm-api"] 