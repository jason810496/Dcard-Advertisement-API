EXECUTABLE := api
SOURCES ?= $(shell find . -name "*.go" -type f)
GO ?= go

build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s \
	-w $(LDFLAGS)' -o bin/$@ ./cmd/$(EXECUTABLE) 

build-generator:
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s \
	-w $(LDFLAGS)' -o bin/generator ./cmd/generator

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
	minikube image load dcard-advertisement-api-generator:latest
	minikube image load dcard-advertisement-api-api:latest
	minikube image load dcard-advertisement-api-scheduler:latest  
	minikube image load dcard-advertisement-api-generator:latest

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

.PHONY: kube-generator
kube-generator:
	kubectl apply -f ./deployments/kubernetes/generator/pod.yaml

.PHONY: kube-del-generator
kube-del-generator:
	kubectl delete -f ./deployments/kubernetes/generator/pod.yaml

.PHONY: kube-all
kube-all: kube-reload-image kube-config kube-database kube-redis kube-api kube-scheduler

.PHONY: kube-del-all
kube-del-all: kube-del-config kube-del-database kube-del-redis kube-del-api kube-del-scheduler

.PHONY: kube-k6-operator
kube-k6-operator:
	kubectl apply -f ./deployments/kubernetes/k6/k6-operator.yaml

.PHONY: kube-del-k6-operator
kube-del-k6-operator:
	kubectl delete -f ./deployments/kubernetes/k6/k6-operator.yaml

.PHONY: kube-k6-config
kube-k6-config:
	kubectl create configmap k6-script-config --from-file=k6/load-test.js
	kubectl create configmap k6-config --from-env-file=.env/dev/k6-operator.env

.PHONY: kube-del-k6-config
kube-del-k6-config:
	kubectl delete configmap k6-script-config
	kubectl delete configmap k6-config

.PHONY: kube-k6-resource
kube-k6-resource:
	kubectl apply -f ./deployments/kubernetes/k6/k6-resource.yaml

.PHONY: kube-del-k6-resource
kube-del-k6-resource:
	kubectl delete -f ./deployments/kubernetes/k6/k6-resource.yaml

.PHONY: build-image-linux
build-image-linux:
	docker build -t dcard-advertisement-api-api:latest -f ./deployments/dev/api/Dockerfile --platform=linux/amd64 . 
	docker build -t dcard-advertisement-api-scheduler:latest -f ./deployments/dev/scheduler/Dockerfile --platform=linux/amd64 .
	docker build -t dcard-advertisement-api-generator:latest -f ./deployments/dev/generator/Dockerfile --platform=linux/amd64 .
	docker build -t dcard-advertisement-api-k6-runner:latest -f ./deployments/dev/k6/Dockerfile.runner --platform=linux/amd64 .
	docker build -t dcard-advertisement-api-k6-starter:latest -f ./deployments/dev/k6/Dockerfile.starter --platform=linux/amd64 .

.PHONY: tag-image
tag-image:
	docker tag dcard-advertisement-api-api:latest jasonbigcow/dcard-advertisement-api-api:latest
	docker tag dcard-advertisement-api-scheduler:latest jasonbigcow/dcard-advertisement-api-scheduler:latest
	docker tag dcard-advertisement-api-generator:latest jasonbigcow/dcard-advertisement-api-generator:latest
	docker tag dcard-advertisement-api-k6-runner:latest jasonbigcow/dcard-advertisement-api-k6-runner:latest
	docker tag dcard-advertisement-api-k6-starter:latest jasonbigcow/dcard-advertisement-api-k6-starter:latest

.PHONY: push-image
push-image:
	docker push jasonbigcow/dcard-advertisement-api-api:latest
	docker push jasonbigcow/dcard-advertisement-api-scheduler:latest
	docker push jasonbigcow/dcard-advertisement-api-generator:latest
	docker push jasonbigcow/dcard-advertisement-api-k6-runner:latest
	docker push jasonbigcow/dcard-advertisement-api-k6-starter:latest

.PHONY: gke-database
gke-database:
	kubectl apply -f ./deployments/gke/database/statefulset.yaml
	kubectl apply -f ./deployments/gke/database/service.yaml

.PHONY: gke-del-database
gke-del-database:
	kubectl delete -f ./deployments/gke/database/statefulset.yaml
	kubectl delete -f ./deployments/gke/database/service.yaml

.PHONY: gke-redis
gke-redis:
	kubectl apply -f ./deployments/gke/redis/statefulset.yaml
	kubectl apply -f ./deployments/gke/redis/service.yaml

.PHONY: gke-del-redis
gke-del-redis:
	kubectl delete -f ./deployments/gke/redis/statefulset.yaml
	kubectl delete -f ./deployments/gke/redis/service.yaml

.PHONY: gke-api
gke-api:
	kubectl apply -f ./deployments/gke/api/deployment.yaml
	kubectl apply -f ./deployments/gke/api/service.yaml

.PHONY: gke-del-api
gke-del-api:
	kubectl delete -f ./deployments/gke/api/deployment.yaml
	kubectl delete -f ./deployments/gke/api/service.yaml

.PHONY: gke-scheduler
gke-scheduler:
	kubectl apply -f ./deployments/gke/scheduler/deployment.yaml

.PHONY: gke-del-scheduler
gke-del-scheduler:
	kubectl delete -f ./deployments/gke/scheduler/deployment.yaml

.PHONY: gke-generator
gke-generator:
	kubectl apply -f ./deployments/gke/generator/pod.yaml

.PHONY: gke-del-generator
gke-del-generator:
	kubectl delete -f ./deployments/gke/generator/pod.yaml

.PHONY: gke-all
gke-all: kube-config gke-database gke-redis gke-api gke-generator

.PHONY: gke-del-all
gke-del-all: gke-del-database gke-del-redis gke-del-api gke-del-generator gke-del-scheduler kube-del-config 

.PHONY: gke-k6-resource
gke-k6-resource:
	kubectl apply -f ./deployments/gke/k6/k6-resource.yaml

.PHONY: gke-del-k6-resource
gke-del-k6-resource:
	kubectl delete -f ./deployments/gke/k6/k6-resource.yaml

.PHONY: gke-k6
gke-k6: kube-k6-operator kube-k6-config gke-k6-resource

.PHONY: gke-del-k6
gke-del-k6: gke-del-k6-resource kube-del-k6-config


.PHONY: grafana-cloud-config
grafana-cloud-config:
	kubectl create configmap grafana-cloud-config --from-env-file=.env/kubernetes/k6-grafana-cloud-prometheus.env

.PHONY: grafana-cloud-del-config
grafana-cloud-del-config:
	kubectl delete configmap grafana-cloud-config

.PHONY: grafana-cloud-secret
grafana-cloud-secret:
	kubectl create secret generic grafana-cloud-secret --from-env-file=.env/kubernetes/k6-grafana-cloud-secret.env

.PHONY: grafana-cloud-del-secret
grafana-cloud-del-secret:
	kubectl delete secret grafana-cloud-secret