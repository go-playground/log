GOCMD=GO111MODULE=on go

linters-install:
	@golangci-lint --version >/dev/null 2>&1 || { \
		echo "installing linting tools..."; \
		brew install golangci-lint; \
	}

lint: linters-install
	@golangci-lint run

test:
	$(GOCMD) test -cover -race ./...

bench:
	$(GOCMD) test -run=NONE -bench=. -benchmem ./...

.PHONY: test lint linters-install