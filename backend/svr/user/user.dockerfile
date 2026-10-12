ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug \
    TZ=Asia/Shanghai

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi
RUN if [ ! -d "/app/configs" ]; then mkdir -p /app/configs; fi

COPY bin/user /app
COPY configs/user_config.yaml /app/configs/user_config.yaml

RUN chmod +x /app/user

ENTRYPOINT ["/app/user"]