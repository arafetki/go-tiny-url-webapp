## Colors
COLOR_RESET   = \033[0m
COLOR_INFO    = \033[32m
COLOR_COMMENT = \033[33m

## Variables
MAIN_PACKAGE_PATH := ./cmd/api
BINARY_NAME := tinyurl

.PHONY: help
## Help
help:
	@printf "${COLOR_COMMENT}Usage:${COLOR_RESET}\n"
	@printf " make [target] [args...]\n\n"
	@printf "${COLOR_COMMENT}Available targets:${COLOR_RESET}\n"
	@awk '/^[a-zA-Z\-\0-9\.@]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf " ${COLOR_INFO}%-16s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: tidy
## format code and tidy modfile
tidy:
	go mod tidy -v
	go fmt ./...


.PHONY: audit
## run quality control checks
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...


.PHONY: test
## run unit tests
test:
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: build
## build the application
build:
	go build -ldflags='-s -w' -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

.PHONY: run
## run the application
run: 
	build
	./bin/${BINARY_NAME}


# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #


.PHONY: migrations/new
## create a new database migration
migrations/new:
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./assets/migrations ${name}


.PHONY: migrations/up
## apply all up database migrations
migrations/up:
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="sqlite3://${DATABASE_DSN}" up


.PHONY: migrations/down
## apply all down database migrations
migrations/down:
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="sqlite3://${DATABASE_DSN}" down


.PHONY: migrations/goto
## migrate to a specific version number
migrations/goto:
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="sqlite3://${DATABASE_DSN}" goto ${version}


.PHONY: migrations/version
## print the current in-use migration version
migrations/version:
	go run -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="sqlite3://${DATABASE_DSN}" version