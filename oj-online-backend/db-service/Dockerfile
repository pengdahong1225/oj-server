FROM golang:1.20-alpine AS budiler

ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn \
    LOG_MODE=release

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir /app/log; fi

COPY . .
COPY ../../config /app/config

RUN go mod tidy
RUN go build -o db_service .

EXPOSE 9000

ENTRYPOINT ["/app/db_service"]