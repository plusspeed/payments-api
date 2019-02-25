
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
base_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
SERVICE = $(base_dir)

BUILDENV += CGO_ENABLED=0

TESTFLAGS := -v -cover
LINT_FLAGS :=--disable-all --enable=vet --enable=vetshadow --enable=golint --enable=ineffassign --enable=goconst --enable=gofmt
LINTER_EXE := gometalinter.v2
LINTER := $(GOPATH)/bin/$(LINTER_EXE)


.PHONY: install
install:
	go get -v -t -d ./...

$(LINTER):
	go get -u gopkg.in/alecthomas/$(LINTER_EXE)
	$(LINTER) --install

.PHONY: lint
lint: $(LINTER)
ifdef LEXC
	$(LINTER) --exclude '$(LEXC)' $(LINT_FLAGS) ./...
else
	$(LINTER) $(LINT_FLAGS) ./...
endif

.PHONY: clean
clean:
	rm -f $(SERVICE)

# build the binary
$(SERVICE):
	$(BUILDENV) go build -o $(SERVICE)

build: clean $(SERVICE)

.PHONY: test
test:
	$(BUILDENV) go test $(TESTFLAGS) ./...

.PHONY: all
all: clean $(LINTER) lint test build

