SOURCE=./cmd/sysmonitor
APP=sysmonitor
GOBASE=$(shell pwd)
RELEASE_DIR=$(GOBASE)/bin
VERSION=1.0

.DEFAULT_GOAL = build 

GO_SRC_DIRS := $(shell \
	find . -name "*.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)
GO_TEST_DIRS := $(shell \
	find . -name "*_test.go" -not -path "./vendor/*" -not -path "./tests/*" | \
	xargs -I {} dirname {}  | \
	uniq)	

build: 
	@CGO_ENABLED=0 go build -v -o sysmonitor ${SOURCE}
	
run:
	go run ${SOURCE} grpc_server --config=configs/prod.yaml --port=50053
	#go run ./cmd/sysmonitor grpc_server --port=50053
	
run-client:
	@go run ${SOURCE} grpc_client --server=":50053"	-t=10
	#go run ./cmd/sysmonitor grpc_client --server=":50051"

lint: 
	@goimports -w ${GO_SRC_DIRS}	
	@gofmt -s -w ${GO_SRC_DIRS}
	@golangci-lint run

test:
	go test ${GO_TEST_DIRS}
	@#go test -race -count 100 ${GO_TEST_DIRS}

gen:
	protoc -I api/proto --go_out=plugins=grpc:pkg/sysmonitor/api api/proto/sysmonitor.proto

mod:
	go mod verify
	go mod tidy

up: build
	docker-compose -f deployments/docker-compose.yml up --build --detach

down:
	docker-compose  --file deployments/docker-compose.yml down

.PHONY: integration-tests
# Запуск Integration tests
integration-tests:
	set -e ;\
	docker-compose -f deployments/docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f deployments/docker-compose.test.yml run integration_tests go test ./tests/integration || test_status_code=$$? ;\
	docker-compose -f deployments/docker-compose.test.yml down --volumes;\
	printf "Return code is $$test_status_code\n" ;\
	exit $$test_status_code ;\

release:
	rm -rf ${RELEASE_DIR}/
	GOOS=windows GOARCH=amd64 go build -o ${RELEASE_DIR}/win/${APP}.exe ${SOURCE}
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ${RELEASE_DIR}/linux/${APP} ${SOURCE}

.PHONY: build run release lint test gen mod up down