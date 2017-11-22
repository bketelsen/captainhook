# modified from https://github.com/cloudflare/hellogopher  
# which is copyrighted by Cloudflare
#
# import path is where your repository can be found.
# To import subpackages, always prepend the full import path.
# If you change this, run `make clean`. Read more: https://git.io/vM7zV
IMPORT_PATH := github.com/bketelsen/captainhook
DOCKER_IMAGE := captainhook
build_dir := $(CURDIR)/build
dist_dir := $(CURDIR)/dist
exec := $(DOCKER_IMAGE)
github_repo := bketelsen/captainhook

# V := 1 # When V is set, print commands and build progress.

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

.PHONY: all
all: test build

.PHONY: build
build: .GOPATH/.ok
	@echo "Building..."
	$Q go install $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)

.PHONY: clean-build
clean-build:
	@echo "Removing build files"
	rm -rf $(build_dir)

.PHONY: clean-dist
clean-dist:
	@echo "Removing distribution files"
	rm -rf $(dist_dir)

.PHONY: tags
tags: .GOPATH/.ok
	@echo "Listing tags..."
	$Q @git tag

.PHONY: release
release: .GOPATH/.ok clean-dist build
	goreleaser


### Code not in the repository root? Another binary? Add to the path like this.
# .PHONY: otherbin
# otherbin: .GOPATH/.ok
#   $Q go install $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/otherbin

##### ^^^^^^ EDIT ABOVE ^^^^^^ #####

##### =====> Utility targets <===== #####

.PHONY: clean test list cover format docker deps

deps: .GOPATH/.ok
	@echo "Ensuring Dependencies..."
	$Q go env
	$Q dep ensure

docker:
	@echo "Docker Build..."
	$Q docker build -t $(DOCKER_IMAGE) .

clean:
	@echo "Clean..."
	$Q rm -rf bin .GOPATH

test: .GOPATH/.ok
	@echo "Testing..."
	$Q go test $(if $V,-v) -i -race $(allpackages) # install -race libs to speed up next run
ifndef CI
	@echo "Testing Outside CI..."
	$Q go vet $(allpackages)
	$Q GODEBUG=cgocheck=2 go test -race $(allpackages)
else
	@echo "Testing in CI..."
	$Q ( go vet $(allpackages); echo $$? ) | \
       tee .GOPATH/test/vet.txt | sed '$$ d'; exit $$(tail -1 .GOPATH/test/vet.txt)
	$Q ( GODEBUG=cgocheck=2 go test -v -race $(allpackages); echo $$? ) | \
       tee .GOPATH/test/output.txt | sed '$$ d'; exit $$(tail -1 .GOPATH/test/output.txt)
endif

list: .GOPATH/.ok
	@echo "List..."
	@echo $(allpackages)

cover: bin/gocovmerge .GOPATH/.ok
	@echo "Coverage Report..."
	@echo "NOTE: make cover does not exit 1 on failure, don't use it to check for tests success!"
	$Q rm -f .GOPATH/cover/*.out .GOPATH/cover/all.merged
	$(if $V,@echo "-- go test -coverpkg=./... -coverprofile=.GOPATH/cover/... ./...")
	@for MOD in $(allpackages); do \
        go test -coverpkg=`echo $(allpackages)|tr " " ","` \
            -coverprofile=.GOPATH/cover/unit-`echo $$MOD|tr "/" "_"`.out \
            $$MOD 2>&1 | grep -v "no packages being tested depend on"; \
    done
	$Q ./bin/gocovmerge .GOPATH/cover/*.out > .GOPATH/cover/all.merged
ifndef CI
	@echo "Coverage Report..."
	$Q go tool cover -html .GOPATH/cover/all.merged
else
	@echo "Coverage Report In CI..."
	$Q go tool cover -html .GOPATH/cover/all.merged -o .GOPATH/cover/all.html
endif
	@echo ""
	@echo "=====> Total test coverage: <====="
	@echo ""
	$Q go tool cover -func .GOPATH/cover/all.merged

format: bin/goimports .GOPATH/.ok
	@echo "Formatting..."
	$Q find .GOPATH/src/$(IMPORT_PATH)/ -iname \*.go | grep -v \
        -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) | xargs ./bin/goimports -w

##### =====> Internals <===== #####

.PHONY: setup
setup: clean .GOPATH/.ok
	@echo "Setup..."
	if ! grep "/.GOPATH" .gitignore > /dev/null 2>&1; then \
        echo "/.GOPATH" >> .gitignore; \
    fi
	if ! grep "/bin" .gitignore > /dev/null 2>&1; then \
        echo "/bin" >> .gitignore; \
    fi
	go get -u github.com/golang/dep/cmd/dep
	go get github.com/wadey/gocovmerge
	go get golang.org/x/tools/cmd/goimports
	go get github.com/mitchellh/gox

VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H:%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

# cd into the GOPATH to workaround ./... not following symlinks
_allpackages = $(shell ( cd $(CURDIR)/.GOPATH/src/$(IMPORT_PATH) && \
    GOPATH=$(CURDIR)/.GOPATH go list ./... 2>&1 1>&3 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) 1>&2 ) 3>&1 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)))

# memoize allpackages, so that it's executed only once and only if used
allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)

export GOPATH := $(CURDIR)/.GOPATH
unexport GOBIN

Q := $(if $V,,@)

.GOPATH/.ok: 
	@echo "Checking GOPATH..."
	$Q mkdir -p "$(dir .GOPATH/src/$(IMPORT_PATH))"
	$Q ln -s ../../../.. ".GOPATH/src/$(IMPORT_PATH)"
	$Q mkdir -p .GOPATH/test .GOPATH/cover
	$Q mkdir -p bin
	$Q ln -s ../bin .GOPATH/bin
	$Q touch $@

.PHONY: bin/gocovmerge bin/goimports
bin/gocovmerge: .GOPATH/.ok
	@echo "Checking Coverage Tool Installation..."
	@test -d ./.GOPATH/src/github.com/wadey/gocovmerge || \
        { echo "Vendored gocovmerge not found, try running 'make setup'..."; exit 1; }
	$Q go install github.com/wadey/gocovmerge
bin/goimports: .GOPATH/.ok
	@echo "Checking Import Tool Installation..."
	@test -d ./.GOPATH/src/golang.org/x/tools/cmd/goimports || \
        { echo "Vendored goimports not found, try running 'make setup'..."; exit 1; }
	$Q go install golang.org/x/tools/cmd/goimports

