GO=go
OUT=terraform-provider-cern

TEST?=$$(go list .)
GOFMT_FILES?=$$(find . -name '*.go')

default: build

build:
	$(ENV) $(GO) build -o $(OUT)

fmt:
	gofmt -s -w $(GOFMT_FILES)

debug:
	$(ENV) $(GO) build -gcflags=all="-N -l" -o $(OUT)
