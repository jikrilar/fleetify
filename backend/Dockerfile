FROM golang:1.26.3-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /fleetify ./cmd/server

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /fleetify /app/fleetify
COPY uploads /app/uploads
EXPOSE 8080
CMD ["/app/fleetify"]
