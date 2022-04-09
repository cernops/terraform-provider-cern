BIN          = terraform-provider-cern
GOFMT_FILES ?= $$(find . -name '*.go')
GO_ARGS     ?=

all: build

$(BIN): cern main.go go.mod go.sum
	go build $(GO_ARGS) -o $@


fmt:
	gofmt -s -w $(GOFMT_FILES)

build: $(BIN)

debug: GO_ARGS += -gcflags=all="-N -l"
debug: $(BIN)

.PHONY: all build fmt debug
