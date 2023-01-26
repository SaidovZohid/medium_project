-include .env
.SILENT:
CURRENT_DIR=$(shell pwd)
DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

run:
	go run cmd/main.go
	
print:
	echo $(DB_URL)

composeup:
	docker compose --env-file ./.env.docker up

pull-sub-module:
	git submodule update --init --recursive

update-sub-module:
	git submodule update --remote --merge 

proto-gen:
	rm -rf genproto
	./scripts/gen-proto.sh ${CURRENT_DIR}