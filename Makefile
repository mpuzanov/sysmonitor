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
	find . -name "*_test.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)	

build: 
	@CGO_ENABLED=0 go build -v -o sysmonitor ${SOURCE}
run:
	@go run ${SOURCE} grpc_server --config=configs/prod.yaml
	
run-client:
	@go run ${SOURCE} grpc_client --server="0.0.0.0:50051"	

lint: 
	@goimports -w ${GO_SRC_DIRS}	
	@gofmt -s -w ${GO_SRC_DIRS}
	@#golangci-lint run

test:
	go test -race -count 100 ${GO_TEST_DIRS}

gen:
	protoc -I api/proto --go_out=plugins=grpc:pkg/sysmonitor/api api/proto/sysmonitor.proto

mod:
	go mod verify
	go mod tidy


release:
	rm -rf ${RELEASE_DIR}${APP}*
	GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o ${RELEASE_DIR}/${APP}.exe ${SOURCE}
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o ${RELEASE_DIR}/${APP} ${SOURCE}

.PHONY: build run release lint test gen mod up down