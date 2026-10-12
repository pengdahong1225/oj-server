ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug \
    TZ=Asia/Shanghai

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi
RUN if [ ! -d "/app/configs" ]; then mkdir -p /app/configs; fi

COPY bin/gateway /app
COPY configs/gateway_config.yaml /app/configs/gateway_config.yaml

RUN chmod +x /app/gateway

ENTRYPOINT ["/app/gateway"]