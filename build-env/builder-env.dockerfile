FROM debian:trixie-slim AS builder-env

# 添加gcc编译环境
RUN apt-get update &&  \
    apt-get install -y gcc g++ make cmake

# 输出版本
RUN <<EOF
gcc --version
make --version
EOF

