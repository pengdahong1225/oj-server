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
WORKDIR /root/builder

VOLUME["/root/builder"]

CMD ["/bin/bash"]