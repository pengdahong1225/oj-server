ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug \
    TZ=Asia/Shanghai

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi
RUN if [ ! -d "/app/configs" ]; then mkdir -p /app/configs; fi

COPY bin/judge_linux /app
COPY configs/judge_config.yaml /app/configs/judge_config.yaml

ENTRYPOINT ["/app/judge_linux"]