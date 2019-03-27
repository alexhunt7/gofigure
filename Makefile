PERCENT := %

.PHONY: all clean fmt lint vet ineffassign misspell cyclo proto test bench upx

all: clean fmt lint vet ineffassign misspell cyclo proto gofigure test upx

clean:
	@rm -rf proto/*.go gofigure coverage.out docs
fmt:
	@go fmt ./...

lint: fmt
	@golangci-lint run

docs: proto fmt
	@bash -c 'godoc -http=:6060 &>/dev/null & sleep 1 && wget --quiet -e robots=off -r -np -N -E -p -k http://localhost:6060/pkg/github.com/alexhunt7/gofigure/; mv "localhost:6060" docs; kill $(PERCENT)1'
	@firefox docs/pkg/github.com/alexhunt7/gofigure/index.html

proto: clean
	@protoc -I proto/ proto/*.proto --go_out=plugins=grpc:proto

gofigure: proto fmt
	@CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o gofigure

test: proto fmt
	@go test ./... -coverprofile=coverage.out
	#@go tool cover -html=coverage.out

bench: proto fmt
	@go test -bench=. ./...

upx: gofigure
	@upx -qq gofigure
