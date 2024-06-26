FROM debian:trixie-slim

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    gcc \
    g++ \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir /app/log; fi

COPY ./go-judge_1.8.5_linux_amd64 /app/go-judge_1.8.5_linux_amd64

EXPOSE 5050/tcp 5051/tcp

ENTRYPOINT ["/app/go-judge_1.8.5_linux_amd64"]