
.PHONY: build

deps:
	go get -t ./...

test:
	go test ./...

test-all:
	go clean -testcache
	go test ./...

lint: golint vet goimports

vet:
	ret=0 && test -z "$$(go vet ./... | tee /dev/stderr)" || ret=1 ; exit $$ret

golint:
	ret=0 && test -z "$$(golint . | tee /dev/stderr)" || ret=1 ; exit $$ret

goimports:
	ret=0 && test -z "$$(goimports -l . | tee /dev/stderr)" || ret=1 ; exit $$ret

tools:
	[ -f $(GOPATH)/bin/goimports ] || go get golang.org/x/tools/cmd/goimports
	[ -f $(GOPATH)/bin/golint ] || go get github.com/golang/lint/golint
