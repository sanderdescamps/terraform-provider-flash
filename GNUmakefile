TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PROVIDER_NAME=flash
VERSION=1.1.3
ARCH=amd64
OS=linux
OS_ARCH=${OS}_${ARCH}
PKG_NAME=terraform-provider-${PROVIDER_NAME}
TF_PLUGIN_PATH=~/.terraform.d/plugins/localdomain/provider/${PROVIDER_NAME}/${VERSION}/${OS}_${ARCH}

default: build

.build: 
	GOOS=${OS} GOARCH=${ARCH} go build -o ${PKG_NAME}

.install: 
	mkdir -p ${TF_PLUGIN_PATH}
	cp ${PKG_NAME} ${TF_PLUGIN_PATH}/${PKG_NAME}
	mv ${PKG_NAME} ${PKG_NAME}_${VERSION}_${OS_ARCH}

build: fmtcheck .build
	
install: build .install
	
test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile website website-test

