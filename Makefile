CORES ?= $(shell nproc)

.PHONY: all _all clean proto

all:
	@$(MAKE) _all --no-print-directory -j$(CORES)

_all: clean proto gofigure-client gofigure-server

clean:
	rm -f proto/*.go gofigure-client gofigure-server

proto:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto

gofigure-server: proto
	go build -o gofigure-server server/*.go

gofigure-client: proto
	go build -o gofigure-client client/*.go
