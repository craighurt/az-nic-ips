DOCKER_IMAGE_NAME := az-nic-ips
DOCKER_CONTAINER_NAME := ${DOCKER_IMAGE_NAME}-build-container_$(shell date +"%s")

default: build

build: clean
	docker build --target builder --tag ${DOCKER_IMAGE_NAME} -f Dockerfile .
	docker create --name ${DOCKER_CONTAINER_NAME} ${DOCKER_IMAGE_NAME}
	mkdir -p bin
	docker cp ${DOCKER_CONTAINER_NAME}:/go/bin/azip bin/azip
	-docker rm -f ${DOCKER_CONTAINER_NAME}

image:
	docker build --tag docker4x/${DOCKER_IMAGE_NAME}:latest -f .

clean:
	-docker rm -f ${DOCKER_CONTAINER_NAME}
	-rm bin/azip
