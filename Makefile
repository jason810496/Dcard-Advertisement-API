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

.PHONY: minikube-reload-image
minikube-reload-image:
	minikube image rm dcard-advertisement-api-api:latest
	minikube image rm dcard-advertisement-api-scheduler:latest
	minikube image load dcard-advertisement-api-api:latest
	minikube image load dcard-advertisement-api-scheduler:latest  

.PHONY: kube-config
kube-config:
	kubectl create configmap api-scheduler-config --from-file=.env/kubernetes.yaml
	kubectl create configmap db-config --from-env-file=.env/dev/db.env
	kubectl create configmap redis-config --from-env-file=.env/dev/redis.env

.PHONY: kube-del-config
kube-del-config:
	kubectl delete configmap api-scheduler-config
	kubectl delete configmap db-config
	kubectl delete configmap redis-config

.PHONY: kube-database
kube-database:
	kubectl apply -f ./deployments/kubernetes/database/statefulset.yaml
	kubectl apply -f ./deployments/kubernetes/database/service.yaml

.PHONY: kube-del-database
kube-del-database:
	kubectl delete -f ./deployments/kubernetes/database/statefulset.yaml
	kubectl delete -f ./deployments/kubernetes/database/service.yaml

.PHONY: kube-redis
kube-redis:
	kubectl apply -f ./deployments/kubernetes/redis/statefulset.yaml
	kubectl apply -f ./deployments/kubernetes/redis/service.yaml

.PHONY: kube-del-redis
kube-del-redis:
	kubectl delete -f ./deployments/kubernetes/redis/statefulset.yaml
	kubectl delete -f ./deployments/kubernetes/redis/service.yaml

.PHONY: kube-api
kube-api:
	kubectl apply -f ./deployments/kubernetes/api/deployment.yaml
	kubectl apply -f ./deployments/kubernetes/api/service.yaml

.PHONY: kube-del-api
kube-del-api:
	kubectl delete -f ./deployments/kubernetes/api/deployment.yaml
	kubectl delete -f ./deployments/kubernetes/api/service.yaml

.PHONY: kube-scheduler
kube-scheduler:
	kubectl apply -f ./deployments/kubernetes/scheduler/deployment.yaml

.PHONY: kube-del-scheduler
kube-del-scheduler:
	kubectl delete -f ./deployments/kubernetes/scheduler/deployment.yaml

.PHONY: kube-all
kube-all: kube-reload-image kube-config kube-database kube-redis kube-api kube-scheduler


.PHONY: kube-del-all
kube-del-all: kube-del-config kube-del-database kube-del-redis kube-del-api kube-del-scheduler

