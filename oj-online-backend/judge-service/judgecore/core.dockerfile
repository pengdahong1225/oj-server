FROM debian:trixie-slim AS builder-env

# 添加gcc编译环境
RUN apt-get update &&  \
    apt-get install -y gcc g++ make

# 输出版本
RUN <<EOF
gcc --version
make --version
EOF


FROM builder-env AS builder-core

RUN apt-get install -y libseccomp-dev

# 定义工作区目录
WORKDIR /root/builder

CMD ["/bin/bash"]