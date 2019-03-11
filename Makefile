GOCMD=GO111MODULE=on go

lint:
	go vet ./...

test:
	$(GOCMD) test -cover -race ./...

bench:
	$(GOCMD) test -run=NONE -bench=. -benchmem ./...

.PHONY: test lint