ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi
RUN if [ ! -d "/app/config" ]; then mkdir -p /app/config; fi

COPY ./service /app

ENTRYPOINT ["/app/service"]