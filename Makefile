.PHONY: wire
wire:
	cd cmd/app && wire

.PHONY: run
run:
	sbin/app run

.PHONY: build
build:
	go build -o sbin/app ./cmd/app

.PHONY: app-init
app-init:
	sbin/app init

.PHONY: restart
restart:
	go build -o sbin/app ./cmd/app
	sbin/app run

.PHONY: mcp
mcp:
	cd cmd/mcp && wire
	go build -o sbin/mcp ./cmd/mcp
	sbin/mcp run

DOCKER_COMPOSE_FILES = -f docker-compose-env.yml -f docker-compose.yml

.PHONY: clean
clean:
	rm -rf data/server

.PHONY: up
up:
	docker compose $(DOCKER_COMPOSE_FILES) up -d