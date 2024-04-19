FROM debian:trixie-slim AS builder

# 添加gcc编译环境
RUN apt-get update &&  \
    apt-get install -y gcc g++ make && \
    apt-get install -y libseccomp-devel

# 输出版本
RUN <<EOF
gcc --version
make --version
EOF

# 定义工作区目录
WORKDIR /app

FROM builder AS compiler

COPY . ./app

VOLUME["/app/lib"]

# 编译Makefile
ENTRYPOINT["make clean;make"]
