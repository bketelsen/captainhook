GO_BIN ?= go

install:
	$(GO_BIN) install -v .
	make tidy

tidy:
	$(GO_BIN) mod tidy

deps:
	$(GO_BIN) get -t ./...
	make tidy

build:
	$(GO_BIN) build -v .
	make tidy

test:
	$(GO_BIN) test ./...
	make tidy

ci-deps:
	$(GO_BIN) get -t ./...

ci-test:
	$(GO_BIN) test -race ./...

lint:
	gometalinter --vendor ./... --deadline=1m --skip=internal
	make tidy

update:
	$(GO_BIN) get -u 
	make tidy
	make test
	make install
	make tidy

release-test:
	$(GO_BIN) test -race ./...
	make tidy

release:
	make tidy
	release --skip-packr -f version.go
	make tidy
