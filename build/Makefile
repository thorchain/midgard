# Image Hashes
CHAIN_SERVICE_IMAGE_HASH=$(shell  grep "Successfully built" ./log/chain-service-build.log | cut -d ' ' -f 3)

# Turn off cache
ifeq (${NOCACHE}, TRUE)
CACHE=--no-cache
endif

config:
	@echo CACHE: ${CACHE}
	@echo ECS_CLUSTER: ${ECS_CLUSTER}
	@echo INFLUX_HOST: ${INFLUX_HOST}
	@echo REGISTRY_ID: ${REGISTRY_ID}
	@echo REPOSITORY: ${REPOSITORY}
	@echo ECR_CHAIN_SVC: ${ECR_CHAIN_SVC}
	@echo CHAIN_SERVICE_IMAGE_HASH: ${CHAIN_SERVICE_IMAGE_HASH}

.PHONY: clean
clean:
	rm -rf ./log

.PHONY: login
login:
	aws configure set aws_access_key_id ${AWS_ACCESS_KEY_ID}
	aws configure set aws_secret_access_key ${AWS_SECRET_ACCESS_KEY}
	aws configure set region ${AWS_REGION}
	$(shell aws ecr get-login --no-include-email --region ap-southeast-2)

.PHONY: logger
logger: login
	mkdir -p ./log

# Build and Deploy statechain then chain-svc
deploy: config deploy-chain-svc

# --------------------------------------------------------------- chain service ------------------------------------------ #

# Step 1
build-chain-svc: logger
	docker build -t ${REGISTRY_ID}.${REPOSITORY}/${ECR_CHAIN_SVC}:latest --build-arg influx_host=${INFLUX_HOST} --build-arg chain_host=${CHAIN_HOST} . | tee ./log/chain-service-build.log

# Step 2
tag-chain-svc: build-chain-svc
	docker tag ${REGISTRY_ID}.${REPOSITORY}/${ECR_CHAIN_SVC}:latest ${REGISTRY_ID}.${REPOSITORY}/${ECR_CHAIN_SVC}:${CHAIN_SERVICE_IMAGE_HASH}

# Step 3
push-chain-svc: tag-chain-svc
	docker push ${REGISTRY_ID}.${REPOSITORY}/${ECR_CHAIN_SVC}:${CHAIN_SERVICE_IMAGE_HASH}
	docker push ${REGISTRY_ID}.${REPOSITORY}/${ECR_CHAIN_SVC}:latest

# Step 4
deploy-chain-svc: push-chain-svc
	aws ecs update-service --force-new-deployment --cluster ${ECS_CLUSTER} --service ${ECR_CHAIN_SVC}
