.PHONY: dep test build

build:
	go build .


dep:
	go mod tidy
	go mod vendor

test:
	go test -mod=vendor -gcflags=all=-l $(shell go list ./... | grep -v mock | grep -v docs) -covermode=count -coverprofile .coverage.cov

fmt:
	golines ./ -m 120 -w --base-formatter gofmt --no-reformat-tags
	gofumpt -l -w .
