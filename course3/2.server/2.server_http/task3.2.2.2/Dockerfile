# Build stage
FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN go mod init task3.2.2.2
RUN go mod tidy
RUN go build main.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["/app/main"]

