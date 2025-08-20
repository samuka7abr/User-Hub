FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o userhub ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/userhub .

ENV APP_SECRET="change-me-32+chars"
ENV APP_PEPPER="dev-pepper"

EXPOSE 8080

CMD ["./userhub"]
