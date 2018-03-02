IMPORT_PATH := github.com/bketelsen/captainhook
DOCKER_IMAGE := captainhook
dist_dir := $(CURDIR)/dist
exec := $(DOCKER_IMAGE)
github_repo := bketelsen/captainhook
GITVERSION ?= dev

# V := 1 # When V is set, print commands and build progress.

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

.PHONY: all
all: test build

.PHONY: build
build:
	@echo "Building..."
	$Q go install $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)

.PHONY: clean-dist
clean-dist:
	@echo "Removing distribution files"
	rm -rf $(dist_dir)

.PHONY: tags
tags:
	@echo "Listing tags..."
	$Q @git tag

echo:
	@echo "MESSAGE " $(MESSAGE)

tag:
	@echo "Creating tag" $(GITVERSION)
	$Q @git tag -a v$(GITVERSION) -m $(GITVERSION)

.PHONY: release
release: clean-dist build tag 
	$Q goreleaser

### Code not in the repository root? Another binary? Add to the path like this.
# .PHONY: otherbin
# otherbin: .GOPATH/.ok
#   $Q go install $(if $V,-v) $(VERSION_FLAGS) $(IMPORT_PATH)/cmd/otherbin

##### ^^^^^^ EDIT ABOVE ^^^^^^ #####

##### =====> Utility targets <===== #####

.PHONY: clean test list cover format docker


docker:
	@echo "Docker Build..."
	$Q docker build -t $(DOCKER_IMAGE):$(VERSION) .

clean:
	@echo "Clean..."
	$Q rm -rf bin
	$Q rm -rf captainhook

test:
	@echo "Testing..."
	$Q go test $(if $V,-v) -i -race $(allpackages) # install -race libs to speed up next run
ifndef CI
	@echo "Testing Outside CI..."
	$Q go vet $(allpackages)
	$Q GODEBUG=cgocheck=2 go test -race $(allpackages)
else
	@echo "Testing in CI..."
	$Q ( go vet $(allpackages); echo $$? ) | \
       tee test/vet.txt | sed '$$ d'; exit $$(tail -1 test/vet.txt)
	$Q ( GODEBUG=cgocheck=2 go test -v -race $(allpackages); echo $$? ) | \
       tee test/output.txt | sed '$$ d'; exit $$(tail -1 test/output.txt)
endif

list:
	@echo "List..."
	@echo $(allpackages)

cover: $(GOPATH)/bin/gocovmerge
	@echo "Coverage Report..."
	@echo "NOTE: make cover does not exit 1 on failure, don't use it to check for tests success!"
	$Q rm -f cover/*.out cover/all.merged
	$(if $V,@echo "-- go test -coverpkg=./... -coverprofile=.GOPATH/cover/... ./...")
	@for MOD in $(allpackages); do \
        go test -coverpkg=`echo $(allpackages)|tr " " ","` \
            -coverprofile=cover/unit-`echo $$MOD|tr "/" "_"`.out \
            $$MOD 2>&1 | grep -v "no packages being tested depend on"; \
    done
	$Q gocovmerge cover/*.out > cover/all.merged
ifndef CI
	@echo "Coverage Report..."
	$Q go tool cover -html cover/all.merged
else
	@echo "Coverage Report In CI..."
	$Q go tool cover -html cover/all.merged -o cover/all.html
endif
	@echo ""
	@echo "=====> Total test coverage: <====="
	@echo ""
	$Q go tool cover -func cover/all.merged

format: bin/goimports
	@echo "Formatting..."
	$Q find . -iname \*.go | grep -v \
        -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) | xargs ./bin/goimports -w

##### =====> Internals <===== #####

.PHONY: setup
setup: clean
	@echo "Setup..."
	if ! grep "dist" .gitignore > /dev/null 2>&1; then \
        echo "dist" >> .gitignore; \
    fi
	if ! grep "cover" .gitignore > /dev/null 2>&1; then \
        echo "cover" >> .gitignore; \
    fi
	mkdir -p cover
	go get -u golang.org/x/vgo
	go get github.com/wadey/gocovmerge
	go get golang.org/x/tools/cmd/goimports

VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H:%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

_allpackages = $(shell ( \
    vgo list ./... 2>&1 1>&3 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) 1>&2 ) 3>&1 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)))

# memoize allpackages, so that it's executed only once and only if used
allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)

unexport GOBIN

Q := $(if $V,,@)


.PHONY: ${GOPATH}/bin/gocovmerge ${GOPATH}/bin/goimports
bin/gocovmerge:
	@echo "Checking Coverage Tool Installation..."
	@test -d $(GOPATH)/bin/gocovmerge || \
	$Q go install github.com/wadey/gocovmerge
bin/goimports:
	@echo "Checking Import Tool Installation..."
	@test -d $(GOPATH)/bin/goimports || \
	$Q go install golang.org/x/tools/cmd/goimports

