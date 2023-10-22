## help: print this help message
help:
	@echo 'Usage: '
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

## run/api: run the cmd/api application
api/run:
	go run ./cmd/api

## db/psql: connect to the db using psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/up: apply all up db migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

## db/migrations/new name=$1: create a new db migration
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}
