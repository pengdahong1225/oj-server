FROM debian:latest

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai \
    PYTHON_VERSION=3.11.8 \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    PATH="/usr/local/bin:$PATH"

# 安装python环境
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        wget \
        build-essential \
        zlib1g-dev \
        libncurses5-dev \
        libgdbm-dev \
        libnss3-dev \
        libssl-dev \
        libsqlite3-dev \
        libreadline-dev \
        libffi-dev \
        && \
    # 下载并编译 Python
    wget https://www.python.org/ftp/python/${PYTHON_VERSION}/Python-${PYTHON_VERSION}.tgz && \
    tar -xzf Python-${PYTHON_VERSION}.tgz && \
    cd Python-${PYTHON_VERSION} && \
    ./configure --enable-optimizations && \
    make -j $(nproc) && \
    make altinstall && \
    # 清理构建文件
    cd .. && \
    rm -rf Python-${PYTHON_VERSION} Python-${PYTHON_VERSION}.tgz && \
    # 安装 pip
    wget https://bootstrap.pypa.io/get-pip.py && \
    python${PYTHON_VERSION%.*} get-pip.py && \
    rm get-pip.py && \
    # 清理系统缓存
    apt-get purge -y --auto-remove wget build-essential && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 验证安装
RUN python${PYTHON_VERSION%.*} --version && \
    pip --version

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir /app/log; fi

COPY go-judge_1.8.5_linux_amd64 /app/go-judge_1.8.5_linux_amd64

EXPOSE 5050/tcp 5051/tcp

ENTRYPOINT ["/app/go-judge_1.8.5_linux_amd64"]