.PHONY: generate
generate:
	go generate ./...

.PHONY: cov-html
cov-html:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: update-types
update-types:
	TYPES_UPDATE=1 go test -cover ./...

.PHONY: examples
examples: examples/*
	FILES="$^" go run compile_examples.go
