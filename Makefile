DC = docker compose
DC_CONFIG = deployment/local.compose.yaml

up:
	$(DC) -f $(DC_CONFIG) up -d

.PHONY: down
down:
	$(DC) -f $(DC_CONFIG) down

.PHONY: restart
restart:
	$(DC) -f $(DC_CONFIG) restart

.PHONY: build
build:
	$(DC) -f $(DC_CONFIG) build

.PHONY: stop
stop:
	$(DC) -f $(DC_CONFIG) stop
