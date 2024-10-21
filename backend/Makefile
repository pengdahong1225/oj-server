define build_docker_image
	${info Building ${2}}
	pwd
	docker build -f ${1} -t ${2} .;
endef

db-service:
	$(call build_docker_image, app/db-service/Dockerfile, "db-service")
judge-service:
	$(call build_docker_image, app/judge-service/Dockerfile, "judge-service")
question-service:
	$(call build_docker_image, app/question-service/Dockerfile, "question-service")

docker-all: db-service judge-service question-service