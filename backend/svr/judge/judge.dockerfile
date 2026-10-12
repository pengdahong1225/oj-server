ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug \
    TZ=Asia/Shanghai

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi
RUN if [ ! -d "/app/configs" ]; then mkdir -p /app/configs; fi

COPY bin/judge /app
COPY configs/judge_config.yaml /app/configs/judge_config.yaml

RUN chmod +x /app/judge

ENTRYPOINT ["/app/judge"]