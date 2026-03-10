.PHONY: run

run:
	swag init -g cmd/api/main.go
	go run ./cmd/api/main.go
	
COVERAGE_EXCLUDE = mocks|main.go|docs|test

test:
	- go test ./... -coverprofile=cover.out
	grep -v "${COVERAGE_EXCLUDE}" cover.out > cover.tmp && mv cover.tmp cover.out
	go tool cover -html=cover.out -o cover.html
