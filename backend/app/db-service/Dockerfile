ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi

COPY ./service /app

EXPOSE 9000
ENTRYPOINT ["/app/service"]