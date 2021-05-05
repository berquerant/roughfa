test:
	@mkdir -p tmp
	@go test ./... -race -cover

test-short:
	@go test ./... -race -cover --short
