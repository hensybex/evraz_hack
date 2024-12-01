# Makefile

.PHONY: build up down logs dev dev-down dev-logs setup-env

build:
	docker-compose -f docker-compose.yaml build

up:
	docker-compose -f docker-compose.yaml up -d

down:
	docker-compose -f docker-compose.yaml down -v

logs:
	docker-compose -f docker-compose.yaml logs -f

dev:
	docker-compose -f docker-compose.dev.yaml up --build

dev-down:
	docker-compose -f docker-compose.dev.yaml down -v

dev-logs:
	docker-compose -f docker-compose.dev.yaml logs -f > dev-logs.txt

setup-env:
	cp .env.example .env
