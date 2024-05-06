FROM registry.cn-chengdu.aliyuncs.com/aliyun_pdh_space/builder:v1.0.0 AS builder

WORKDIR /app/judge-service

RUN if [ ! -d "/app/log" ]; then mkdir /app/log; fi

COPY . .
COPY ../../config /config

# 编译
RUN make -f Makefile

EXPOSE 9020

ENTRYPOINT ["/app/judge-service/JudgeService"]