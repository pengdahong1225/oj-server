ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug \
    TZ=Asia/Shanghai

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi
RUN if [ ! -d "/app/config" ]; then mkdir -p /app/config; fi

COPY bin/problem_linux /app
COPY config/problem_config.yaml /app/config/problem_config.yaml
COPY config/app_config.yaml /app/config/app_config.yaml

ENTRYPOINT ["/app/problem_linux"]