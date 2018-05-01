# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=exec
BINARY_UNIX=$(BINARY_NAME)_unix

all: run

build:
	$(GOBUILD) -o ./gateway/build/$(BINARY_NAME) -v ./gateway/cmd/gateway
	$(GOBUILD) -o ./services/order/build/$(BINARY_NAME) -v ./services/order/cmd/order
	$(GOBUILD) -o ./services/ucenter/build/$(BINARY_NAME) -v ./services/ucenter/cmd/ucenter

test:
	$(GOTEST) -v ./

clean:
	$(GOCLEAN)
	rm -f ./gateway/build/$(BINARY_NAME)
	rm -f ./gateway/build/$(BINARY_UNIX)
	rm -f ./services/order/build/$(BINARY_NAME)
	rm -f ./services/order/build/$(BINARY_UNIX)
	rm -f ./services/ucenter/build/$(BINARY_NAME)
	rm -f ./services/ucenter/build/$(BINARY_UNIX)

run:
	$(GOBUILD) -o ./gateway/build/$(BINARY_NAME) -v ./gateway/cmd/gateway
	./gateway/build/$(BINARY_NAME) &
	$(GOBUILD) -o ./services/order/build/$(BINARY_NAME) -v ./services/order/cmd/order
	./services/order/build/$(BINARY_NAME) &
	$(GOBUILD) -o ./services/ucenter/build/$(BINARY_NAME) -v ./services/ucenter/cmd/ucenter
	./services/ucenter/build/$(BINARY_NAME) &

restart:
	kill -INT $$(cat pid)
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v ./
	./build/$(BINARY_NAME)

deps:
	$(GOGET) github.com/kardianos/govendor
	cd ./gateway && govendor sync
	cd ./services/order && govendor sync
	cd ./services/ucenter && govendor sync

cross:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./gateway/build/$(BINARY_NAME) -v ./gateway/cmd/gateway
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./services/order/build/$(BINARY_NAME) -v ./services/order/cmd/order
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./services/ucenter/build/$(BINARY_NAME) -v ./services/ucenter/cmd/ucenter