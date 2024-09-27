ARG image="serverbase"
FROM "${image}"

ENV LOG_MODE=debug

WORKDIR /app

COPY ./service /app

EXPOSE 9000
ENTRYPOINT ["/app/service"]