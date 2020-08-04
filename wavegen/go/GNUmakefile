include ../opinionated.mk

test:
> go test ./...

fmt:
> go fmt ./...
.PHONY: fmt

lint:
> golint ./...
.PHONY: lint

coverage:
> go test -timeout 30s -coverprofile /dev/null ./...
.PHONY: coverage

viewcoverage:
> go test -timeout 30s -coverprofile cover.out ./...
> go tool cover -html=cover.out
.PHONY: viewcoverage

clean:
> rm -rf cover.out
.PHONY: clean
