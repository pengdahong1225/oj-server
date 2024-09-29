services=("question-service" "judge-service" "db-service")

# 退休
for s in "${services[@]}"; do
    docker kill -s 10 "$s"
done

docker-compose -f docker-compose.service.yml stop
