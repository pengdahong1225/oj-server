FROM builder AS runner

ENV LOG_MODE=release

WORKDIR /app

RUN if [ ! -d "/app/log" ]; then mkdir /app/log; fi

COPY ./bin/judge_service /app/judge_service

EXPOSE 8020

ENTRYPOINT ["/app/judge_service"]