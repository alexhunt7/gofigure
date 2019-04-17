PERCENT := %

.PHONY: all clean fmt lint proto test bench upx

all: clean fmt lint proto gofigure test upx

clean:
	@rm -rf proto/*.go gofigure cov/* docs
fmt:
	@go fmt ./...

lint: fmt
	@golangci-lint run

docs: proto fmt
	@bash -c 'godoc -http=:6060 &>/dev/null & sleep 1 && wget --quiet -e robots=off -r -np -N -E -p -k http://localhost:6060/pkg/github.com/alexhunt7/gofigure/; mv "localhost:6060" docs; kill $(PERCENT)1'
	@firefox docs/pkg/github.com/alexhunt7/gofigure/index.html

proto: clean
	@protoc -I proto/ --go_out=plugins=grpc,paths=source_relative:proto proto/*.proto

gofigure: proto fmt
	@CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o gofigure

test: proto fmt
	@go test -coverprofile=cov/coverage.out --coverpkg="$$(go list ./... | grep -v /proto | paste -sd, -)" ./...
	@go tool cover -func=cov/coverage.out
	@go tool cover -html=cov/coverage.out -o cov/coverage.html

bench: proto fmt
	@go test -bench=. ./...

upx: gofigure
	@upx -qq gofigure
