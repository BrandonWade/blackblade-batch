FROM golang:1.15-alpine AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o batch

FROM golang:alpine
WORKDIR /app
COPY --from=builder /app/batch .

CMD ["./batch"]
