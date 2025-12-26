.PHONY: generate-api
generate-api:
	@./scripts/generate.sh

.PHONY: cov
cov:
	go test -cover ./internal/services/profile_service/...

.PHONY: cov-func
cov-func:
	go test -coverprofile=coverage.out ./internal/services/profile_service/...
	go tool cover -func=coverage.out

.PHONY: cov-html
cov-html:
	go test -coverprofile=coverage.out ./internal/services/profile_service/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

.PHONY: mock
mock:
	mockery

.PHONY: run
run:
	@set configPath=./config.yaml && set swaggerPath=./internal/pb/swagger/profile_management_api/profile_management.swagger.json && go run cmd/app/main.go
