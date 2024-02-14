EXECUTABLE := api
SOURCES ?= $(shell find . -name "*.go" -type f)
GO ?= go

build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s \
	-w $(LDFLAGS)' -o bin/$@ ./cmd/$(EXECUTABLE) \

init:
# tools
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/cosmtrek/air@latest
# install dependencies
	$(GO) mod download

.PHONY: gen
gen:
	swag init --parseDependency --parseInternal --parseDepth 2 -g cmd/api/main.go

.PHONY: air
air:
	@hash air > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install github.com/cosmtrek/air@latest; \
	fi

.PHONY: fmt
fmt:
	$(GO) fmt ./...
	swag fmt
	
# run air
.PHONY: dev
dev: air
	air --build.cmd "make" \
	--build.bin bin/api \
	--build.pre_cmd "make gen" \
	--build.exclude_dir "docs" 
