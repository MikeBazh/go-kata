FROM golang:1.19-alpine as builder

WORKDIR go-kata/course3/2.server/2.server_http/task3.2.2.2/main.go

COPY . .

RUN go mod init task3.2.2.2

RUN go mod tidy

RUN go build main.go

CMD ["./main"]

