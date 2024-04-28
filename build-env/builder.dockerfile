FROM registry.cn-chengdu.aliyuncs.com/aliyun_pdh_space/builder-env:v1.0.0 AS builder

# 安装
RUN apt-get install -y libseccomp-dev && \
    apt-get install -y libboost-dev && \
    apt-get install -y protobuf-compiler libprotobuf-dev

# 定义工作区目录
WORKDIR /root/builder

CMD ["/bin/bash"]

