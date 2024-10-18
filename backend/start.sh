services=("question-service" "judge-service" "db-service")

for s in "${services[@]}"; do
    docker stop $s
    docker rm -f $s
done

docker-compose -f docker-compose.service.yml up -d