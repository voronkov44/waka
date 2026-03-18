container_runtime := $(shell which podman || which docker)
$(info using ${container_runtime})

ENV_FILE ?= .env
COMPOSE := ${container_runtime} compose --env-file $(ENV_FILE)

BACKUP_DIR ?= backups
BACKUP_NAME ?= waka_$(shell date +%Y-%m-%d_%H-%M-%S).dump
DB_SERVICE ?= postgres
DB_USER ?= postgres
DB_NAME ?= postgres

.PHONY: up down clean logs ps migrate-docker test test-tasks lint backup restore

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

backup:
	mkdir -p $(BACKUP_DIR)
	$(COMPOSE) exec -T $(DB_SERVICE) pg_dump -U $(DB_USER) -d $(DB_NAME) -Fc > $(BACKUP_DIR)/$(BACKUP_NAME)
	@echo "Backup created: $(BACKUP_DIR)/$(BACKUP_NAME)"

restore:
	@test -n "$(FILE)" || (echo "Usage: make restore FILE=backups/your.dump" && exit 1)
	cat $(FILE) | $(COMPOSE) exec -T $(DB_SERVICE) pg_restore -U $(DB_USER) -d $(DB_NAME) --clean --if-exists
	@echo "Restore completed from: $(FILE)"