.PHONY: generate
generate:
	go generate ./...

.PHONY: cov-html
cov-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: test
test:
	go test -cover ./...

.PHONY: update-types
update-types:
	UPDATE_TYPES=1 go test -cover ./...
