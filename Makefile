.PHONY: build
build:
	goreleaser release

.PHONY: test
test:
	go test ./...
