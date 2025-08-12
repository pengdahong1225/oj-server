FROM debian:latest

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai \
    GOPATH=/go \
    GOROOT=/usr/local/go \
    PATH=$PATH:/usr/local/go/bin

# 安装go环境
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        wget \
        git \
        build-essential && \
    # 下载并安装 Go
    wget -q https://go.dev/dl/go1.23.10.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.10.linux-amd64.tar.gz && \
    rm go1.21.10.linux-amd64.tar.gz && \
    # 清理缓存
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 验证安装
RUN go version && \
    go env

WORKDIR /app

RUN mkdir -p /app/log

COPY go-judge_1.8.5_linux_amd64 /app/go-judge_1.8.5_linux_amd64

EXPOSE 5050/tcp 5051/tcp

ENTRYPOINT ["/app/go-judge_1.8.5_linux_amd64"]