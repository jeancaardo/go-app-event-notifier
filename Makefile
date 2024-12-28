.PHONY: install
export GOSUMDB=off
install:
	@echo "=> Install service"
	@git config --local core.hooksPath .githooks/
	@npm i
	@git config --global url."ssh://git@gitlab.com/".insteadOf https://gitlab.com/
	@git config --global url."ssh://git@github.com/".insteadOf https://github.com/
	@go mod tidy
	@go mod download

.PHONY: build
build:
	@echo "=> Building service"
	@git config --local core.hooksPath .githooks/

	@echo "=> Zipping binaries"

.PHONY: format
format:
	@go fmt ./internal/... ./pkg/... ./cmd/...

.PHONY: start
start:
	@echo "=> Starting service"
	@docker compose up -d
	@make format
	@make build
	@npx sls offline local-authorizers

.PHONY: test
test:
	@echo "=> Running tests"
	@${GOPATH}/bin/golangci-lint run ./internal/... ./pkg/...
	@go test ./internal/... ./pkg/... -covermode=atomic -coverpkg=./... -count=1 -race

.PHONY: test-cover
test-cover:
	@echo "=> Running tests and generating report"
	@go test ./internal/... ./pkg/... -covermode=atomic -coverpkg=./... -cover -coverprofile=c.out -count=1 -race
	@go tool cover -html=c.out -o coverage.html

.PHONY: deploy
deploy:
	@echo "=> Running deployment"
	@make build
	@sls deploy