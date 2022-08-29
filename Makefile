
lint:
	echo "Running staticcheck"
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

vet:
	echo "running go vet"
	go vet

fmt:
	go fmt

code-check: lint vet fmt

unit_test:
	echo "Running unit tests"
	go test -v ./clients

integration_test:
	echo "Running integration test"
	go test -v ./ome

test: unit_test integration_test
	
generate:
	go generate

download:
	go mod download

build: download
	mkdir -p out
	go build -v -o ./out

all: download code-check test compile
	
