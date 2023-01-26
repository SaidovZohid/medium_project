-include .env
.SILENT:
DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable
CURRENT_DIR=$(shell pwd)

run:
	go run cmd/main.go
	
print:
	echo $(DB_URL)

swag-init:
	swag init -g api/api.go -o api/docs

composeup:
	docker compose --env-file ./.env.docker up

proto-gen:
	rm -rf genproto
	./scripts/gen-proto.sh ${CURRENT_DIR}

update-sub-module:
	git submodule update --remote --merge

pull-sub-module:
	git submodule update --init --recursive

test:
	go test -v -cover ./api/v1/...

cache:
	go clean -testcache

generate:
	go generate ./...