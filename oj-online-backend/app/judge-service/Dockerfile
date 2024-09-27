FROM golang:1.21-alpine AS budiler

ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn \
    GIN_MODE=debug \
    LOG_MODE=debug

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir /app/log; fi

COPY . .
COPY ../../config /config

RUN go mod tidy
RUN go build -o judge_service .

EXPOSE 9010

ENTRYPOINT ["/app/judge_service"]