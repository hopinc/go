.PHONY: generate
generate:
	go generate ./...

.PHONY: cov-html
cov-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
