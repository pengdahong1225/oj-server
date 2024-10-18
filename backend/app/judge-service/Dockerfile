ARG image="debian:latest"
FROM "${image}"

ENV LOG_MODE=debug

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir -p /app/log; fi

COPY ./service /app
COPY ./bodyConfig /app

EXPOSE 9010
ENTRYPOINT ["/app/service"]