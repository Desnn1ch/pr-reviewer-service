FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/app

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .
COPY migrations ./migrations
COPY ../temp/pr-reviewer-service/openapi.yml ./openapi.yml

COPY internal/config/config.yaml ./internal/config/config.yaml

ENV APP_PORT=8080
EXPOSE 8080

CMD ["./server"]
