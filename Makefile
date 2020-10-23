.PHONY: test build

build:
	docker build --force-rm=true -t dot/grabber -f build/grabber.Dockerfile .

init:
	docker swarm init --advertise-addr 127.0.0.1 2>/dev/null || true
	docker network create --attachable --driver overlay grabber_net 2>/dev/null || true

run-db:
	docker stack deploy --compose-file deployments/db.yml grabber

run:
	docker stack deploy --compose-file deployments/db.yml grabber
	docker stack deploy --compose-file deployments/backend.yml grabber

stop:
	docker stack rm grabber

local-test: run-db
	go test -p=1 -count=1 ./...

SWAGGER_IMAGE=quay.io/goswagger/swagger:v0.25.0

gen-grabber-server:
	docker run --rm -v $(PWD):/grabber -w /grabber -t $(SWAGGER_IMAGE) \
		generate server \
		--target=internal/rest \
		--exclude-main \
		-f api/grabber.swagger.yml

gen-grabber-client:
	docker run --rm -v $(PWD):/grabber -w /grabber -t $(SWAGGER_IMAGE) \
		generate client \
		--target=internal/rest \
		-f api/grabber.swagger.yml
