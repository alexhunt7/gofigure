CORES ?= $(shell nproc)

.PHONY: all _all clean proto test

all: clean proto gofigure test

clean:
	rm -f proto/*.go gofigure coverage.out

proto:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto

gofigure: proto
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o gofigure

test: proto
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
