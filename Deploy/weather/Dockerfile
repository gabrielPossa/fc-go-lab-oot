# Builder stage
FROM golang:1.26.4-alpine3.23 AS builder

WORKDIR /builder/app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -o server ./cmd/weather

# Final stage
FROM alpine:3.23
WORKDIR /opt/weather

COPY --from=builder /builder/app/server .

EXPOSE 8080

CMD ["./server"]
