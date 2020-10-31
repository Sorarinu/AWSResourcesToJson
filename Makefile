SRCS            := $(shell find ./src -type f -name '*.go' | grep -v vendor)
GO_IMAGE_ID_CMD := docker-compose ps -q go

.PHONY: go/*

default: go/build

go/build: $(SRCS) src/go.mod
	docker-compose up -d --build go
	docker cp `$(GO_IMAGE_ID_CMD)`:/app/main.handle ./bin/main.handle

go/fmt:
	docker-compose run --rm -v $(PWD)/src:/app go make fmt FLAGS="-l -w"

go/lint:
	docker-compose run --rm -v $(PWD)/src:/app go make --no-print-directory -k lint

go/test:
	docker-compose run --rm -v $(PWD)/src:/app go make test