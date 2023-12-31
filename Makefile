# include env vars from .envrc
include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage: '
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
api/run:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN} -smtp-username=${GREENLIGHT_SMTP_USR} -smtp-password=${GREENLIGHT_SMTP_PASS} -smtp-sender=${GREENLIGHT_SMTP_SENDER}

## db/psql: connect to the db using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/up: apply all up db migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

## db/migrations/new name=$1: create a new db migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## format: just format all code
.PHONY: format
format:
	go fmt ./...

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module deps...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module deps...'
	go mod tidy
	go mod verify
	@echo 'Vendoring deps...'
	go mod vendor


# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## build/run/api: run built version of cmd/api
.PHONY: build/run/api
build/run/api:
	./bin/api -port=4040 -db-dsn=${GREENLIGHT_DB_DSN} -smtp-username=${GREENLIGHT_SMTP_USR} -smtp-password=${GREENLIGHT_SMTP_PASS} -smtp-sender=${GREENLIGHT_SMTP_SENDER}
