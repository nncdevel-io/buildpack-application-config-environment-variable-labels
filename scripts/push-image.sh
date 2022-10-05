export DOCKER_IMAGE=${1}
docker tag paketo-dd-java-agent ${DOCKER_IMAGE}
docker push ${DOCKER_IMAGE}
