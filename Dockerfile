FROM golang:1.16.2 AS builder

WORKDIR /go/src
ADD . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    make test && \
    make
#RUN go mod tidy
#RUN go mod vendor

FROM alpine:latest
COPY --from=builder /go/src/.build/app-server /usr/local/bin/app-server
ENTRYPOINT ["/usr/local/bin/app-server"]
