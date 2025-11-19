.PHONY: build
build:
	goreleaser release --clean

.PHONY: test
test:
	go test ./...

.PHONY: convert
convert:
	go run cmd/pgn-tools/main.go convert ~/Downloads/SetupMega2020/Bases/MegaDatabase2020/Mega\ Database\ 2020.cbh output.pgn --verbose --experimental
