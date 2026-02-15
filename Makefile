container_runtime := $(shell which podman || which docker)
$(info using ${container_runtime})

ENV_FILE ?= .env
COMPOSE := ${container_runtime} compose --env-file $(ENV_FILE)

.PHONY: up down clean logs ps

# all services
up: down
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

clean:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f --tail=200

ps:
	$(COMPOSE) ps

migrate-docker:
	$(COMPOSE) run --rm --entrypoint /migrate api -config /config.yaml

test:
	go test ./...

test-tasks:
	@#$(MAKE) -C services test SERVICE=tasks


lint:
	#make -C services lint
