
lint:
	echo "Running staticcheck"
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
	golint ./...

vet:
	echo "running go vet"
	go vet

fmt:
	go fmt ./...

code_check: lint vet fmt

unit_test:
	echo "Running unit tests"
	go test -v ./clients -cover -timeout 60m

integration_test:
	echo "Running integration test"
	TF_ACC=1 go test -v ./ome -timeout 5h -cover

test: unit_test integration_test
	
generate:
	go generate

download:
	go mod download

build: download
	mkdir -p out
	go build -v -o ./out

all: download code-check test compile