.PHONY: run test tidy fmt openapi-validate openapi-generate

run:
	go run ./cmd/core-gateway

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w ./cmd ./internal

openapi-validate:
	@test -f api/openapi.yaml
	@grep -q "openapi: 3.0.3" api/openapi.yaml
	@echo "openapi spec looks present"

openapi-generate:
	@sh scripts/generate-openapi-stubs.sh
