FROM golang:1.16.2 AS builder

WORKDIR /go/src
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o app-server

FROM alpine:latest
COPY --from=builder /go/src/app-server /usr/local/bin/app-server
ENTRYPOINT ["/usr/local/bin/app-server"]
