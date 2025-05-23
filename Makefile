TEST?=$$(go list ./... | grep -v 'vendor')
INSTALL_ROOT?=~/.terraform.d/plugins
HOSTNAME=registry.terraform.io
NAMESPACE=dell
NAME=ome
BINARY=terraform-provider-${NAME}
VERSION?=1.2.3
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64


install: build
	rm -rfv $(INSTALL_ROOT)/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	find examples -type d -name ".terraform" -exec rm -rfv "{}" +;
	find examples -type f -name "trace.*" -delete
	find examples -type f -name "*.tfstate" -delete
	find examples -type f -name "*.hcl" -delete
	find examples -type f -name "*.backup" -delete
	rm -rf trace.*
	
	mkdir -p $(INSTALL_ROOT)/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} $(INSTALL_ROOT)/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

uninstall:
	rm -rfv $(INSTALL_ROOT)/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	find examples -type d -name ".terraform" -exec rm -rfv "{}" +;
	find examples -type f -name "trace.*" -delete
	find examples -type f -name "*.tfstate" -delete
	find examples -type f -name "*.hcl" -delete
	find examples -type f -name "*.backup" -delete
	rm -rf trace.*

unit_test:
	echo "Running unit tests"
	go test -v ./clients -cover -timeout 60m
	go test -v ./utils

integration_test:
	echo "Running integration test"
	TF_ACC=1 go test -v ./ome -timeout 5h -cover

combine_test: unit_test integration_test

test: check
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

check:
	gofmt -s -w .
	go vet
	golangci-lint run --fix --timeout 5m

gosec:
	gosec -quiet ./...

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

generate:
	# Do nothing for now because of vulnerability in terraform-plugin-docs
	# terraform fmt -recursive examples/
	# go generate ./...

cover:
	rm -f coverage.*
	go test -coverprofile=coverage.out ./...
	go tool cover -html coverage.out -o coverage.html

sweep :
	go test -v ./ome -timeout 5h -sweep=all

format:
	gofmt -s -w .	