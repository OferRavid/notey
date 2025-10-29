FROM golang:1.23 as builder
WORKDIR /app
COPY ./internal /app/internal
COPY ./sql /app/sql
COPY ./static /app/static
COPY ./go.mod /app/
COPY main.go /app/
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o notey .

# Stage 2: final minimal image
FROM alpine:latest
WORKDIR /

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/notey /notey
COPY --from=builder /app/static /static
CMD ["/notey"]