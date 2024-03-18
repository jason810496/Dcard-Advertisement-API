EXECUTABLE := api
SOURCES ?= $(shell find . -name "*.go" -type f)
GO ?= go

build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s \
	-w $(LDFLAGS)' -o bin/$@ ./cmd/$(EXECUTABLE) 

build-fake-data:
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s \
	-w $(LDFLAGS)' -o bin/fake-data ./cmd/fake-data

build-scheduler:
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s \
	-w $(LDFLAGS)' -o bin/scheduler ./cmd/scheduler

init:
# tools
	@hash swag > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@hash air > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/cosmtrek/air@latest; \
	fi
# install dependencies
	$(GO) mod download

.PHONY: gen
gen:
	swag init --parseDependency --parseInternal --parseDepth 2 -g cmd/api/main.go

.PHONY: fmt
fmt:
	$(GO) fmt ./...
	swag fmt
	
# run hot reload for local development
.PHONY: local-dev
local-dev:
	air --build.cmd "make build" \
	--build.exclude_dir "stateful_volumes,docs,assets,deployments,bin" \
	--build.bin bin/api \
	--build.args_bin "-config local"
	--build.pre_cmd "make gen" 

.PHONY: local-db
local-db:
	docker compose up db -d
	docker compose up redis -d

.PHONY: local-clean
local-clean:
	docker compose down
	rm -r stateful_volumes/*