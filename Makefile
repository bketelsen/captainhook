IMPORT_PATH := github.com/bketelsen/captainhook
DOCKER_IMAGE := captainhook
dist_dir := $(CURDIR)/dist
exec := $(DOCKER_IMAGE)
github_repo := bketelsen/captainhook
GITVERSION ?= dev
BINARY := captainhook

# V := 1 # When V is set, print commands and build progress.

.PHONY: all
all: setup test build shoulders

.PHONY: build
build: setup
	@echo "Building..."
	$Q go build $(if $V,-v) $(VERSION_FLAGS)

.PHONY: dist
dist: setup clean-dist
	@echo "Building Distribution..."
	$Q go build $(if $V,-v) $(VERSION_FLAGS) -o $(dist_dir)/$(BINARY) $(IMPORT_PATH)

.PHONY: clean-dist
clean-dist:
	@echo "Removing distribution files"
	rm -rf $(dist_dir)

.PHONY: tags
tags:
	@echo "Listing tags..."
	$Q @git tag

tag:
	@echo "Creating tag" $(GITVERSION)
	$Q @git tag -a v$(GITVERSION) -m $(GITVERSION)
	@echo "pushing tag" $(GITVERSION)
	$Q @git push origin v$(GITVERSION)

.PHONY: release
release: setup clean-dist build tag
	$Q goreleaser

shoulders:
	$Q shoulders -w


##### ^^^^^^ EDIT ABOVE ^^^^^^ #####

##### =====> Utility targets <===== #####

.PHONY: clean test list format docker


docker:
	@echo "Docker Build..."
	$Q docker build -t $(DOCKER_IMAGE):$(VERSION) .

clean:
	@echo "Clean..."
	$Q rm -rf $(BINARY)

test:
	@echo "Testing..."
	$Q go test $(if $V,-v) ./...
ifndef CI
	@echo "Testing Outside CI..."
	@echo "GO Vet"
	$Q go vet ./...
	@echo "GO test -race"
	$Q GODEBUG=cgocheck=2 go test -race
else
	@echo "Testing in CI..."
	$Q ( go vet ./...; echo $$? ) | \
       tee test/vet.txt | sed '$$ d'; exit $$(tail -1 test/vet.txt)
	$Q ( GODEBUG=cgocheck=2 go test -v -race ./...; echo $$? ) | \
       tee test/output.txt | sed '$$ d'; exit $$(tail -1 test/output.txt)
endif

list:
	@echo "List..."
	go list -m


format: $(GOIMPORTS)
	@echo "Formatting..."
	$Q find . -iname \*.go | grep -v \
        -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) | xargs $(GOPATH)/bin/goimports -w

##### =====> Internals <===== #####

.PHONY: setup
setup: clean
	@echo "Setup..."
	if ! grep "dist" .gitignore > /dev/null 2>&1; then \
        echo "dist" >> .gitignore; \
    fi
	go get -u rsc.io/goversion
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/gobuffalo/shoulders

VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H:%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'


unexport GOBIN

Q := $(if $V,,@)


GOIMPORTS := $(GOPATH)/bin/goimports

$(GOIMPORTS):
	@echo "Checking Import Tool Installation..."
	@test -d $(GOPATH)/bin/goimports || \
	$Q go install golang.org/x/tools/cmd/goimports

