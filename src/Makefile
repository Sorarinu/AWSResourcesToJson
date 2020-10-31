SRCS := $(shell find . -type f -name '*.go' | grep -v vendor)

.PHONY: default build vet fmt goget golint lint test

default: build 

build: $(SRCS) test
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -ldflags '-s -w' -o ./main.handle ./main.go

vet:
	go vet ./...

fmt:
	echo $(SRCS) | xargs -n 1 gofmt $(FLAGS)

goget: 
	go get -u golang.org/x/lint/golint

gomod: go.mod
	go env -w GO111MODULE=on
	go mod download

golint: goget
	golint -min_confidence=0 -set_exit_status ./...

lint: gomod
	@echo '##### gofmt'
	@$(MAKE) fmt FLAGS=-l
	@echo
	@echo '##### go vet'
	@$(MAKE) vet
	@echo
	@echo '##### golint'
	@$(MAKE) golint

test: lint
	go test ./...