.PHONY: dep


dep:
	go mod tidy
	go mod vendor

test:
	go test -mod=vendor -gcflags=all=-l $(shell go list ./... | grep -v mock | grep -v docs) -covermode=count -coverprofile .coverage.cov
