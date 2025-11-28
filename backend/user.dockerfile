ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug \
    TZ=Asia/Shanghai

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi
RUN if [ ! -d "/app/config" ]; then mkdir -p /app/config; fi

COPY bin/user_linux /app
COPY config/user_config.yaml /app/config/user_config.yaml
COPY config/app_config.yaml /app/config/app_config.yaml

ENTRYPOINT ["/app/user_linux"]