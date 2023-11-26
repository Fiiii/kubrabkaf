# ==============================================================================
# local development
docker-build:
	docker build \
		-f zarf/docker/dockerfile \
		-t kubrabkafka \
		--build-arg BUILD_REF=$(VERSION) \
		.

docker-run: docker-build
	docker run -it \
	  -p 8080:8080 \
	  kubrabkafka