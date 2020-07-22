include ./opinionated.mk

test:
> go test -timeout 30s ./go/...
.PHONY: test

fmt:
> go fmt ./go/...
.PHONY: fmt

lint:
> golint ./go/...
> errcheck ./...
> ineffassign ./
.PHONY: lint

coverage:
> go test -timeout 30s -coverprofile /dev/null ./go/...
.PHONY: coverage

viewcoverage:
> go test -timeout 30s -coverprofile cover.out ./go/...
> go tool cover -html=cover.out
.PHONY: viewcoverage

clean:
> rm -rf cover.out
> go clean
> go clean -testcache
.PHONY: clean
