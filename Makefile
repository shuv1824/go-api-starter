run: build
	@./bin/apiserver

build:
	@go build -o bin/apiserver .
