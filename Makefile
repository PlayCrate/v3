build:
	@go1.19 build -o bin/v3

run: build
	@./bin/v3

test:
	@go1.19 test -v ./...