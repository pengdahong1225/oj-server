FROM golang:1.21-alpine AS budiler

ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn \
    GIN_MODE=release \
    LOG_MODE=release

WORKDIR /app/question-service

RUN if [ ! -d "/app/log" ]; then mkdir /app/log; fi

COPY . .
COPY ../../config /config

RUN go mod tidy
RUN go build -o question_service .

EXPOSE 9020

ENTRYPOINT ["/app/question-service/question_service"]