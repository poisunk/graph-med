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