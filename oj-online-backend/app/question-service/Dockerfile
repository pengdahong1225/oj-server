ARG image="serverbase"
FROM "${image}"

ENV LOG_MODE=debug

WORKDIR /app

COPY ./service /app

EXPOSE 9020
ENTRYPOINT ["/app/service"]