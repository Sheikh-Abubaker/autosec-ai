FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server .

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/server .

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]