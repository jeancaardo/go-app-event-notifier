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
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/users-store/bootstrap					cmd/users/store/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/users-get/bootstrap						cmd/users/get/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/users-getall/bootstrap					cmd/users/get-all/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/users-update/bootstrap					cmd/users/update/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/users-delete/bootstrap					cmd/users/delete/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/events-store/bootstrap					cmd/events/store/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/events-get/bootstrap						cmd/events/get/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/events-getall/bootstrap					cmd/events/get-all/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/events-update/bootstrap					cmd/events/update/main.go
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/events-delete/bootstrap					cmd/events/delete/main.go

	@echo "=> Zipping binaries"
	zip -j bin/users-store.zip bin/users-store/bootstrap
	zip -j bin/users-get.zip bin/users-get/bootstrap
	zip -j bin/users-getall.zip bin/users-getall/bootstrap
	zip -j bin/users-update.zip bin/users-update/bootstrap
	zip -j bin/users-delete.zip bin/users-delete/bootstrap
	zip -j bin/events-store.zip bin/events-store/bootstrap
	zip -j bin/events-get.zip bin/events-get/bootstrap
	zip -j bin/events-getall.zip bin/events-getall/bootstrap
	zip -j bin/events-update.zip bin/events-update/bootstrap
	zip -j bin/events-delete.zip bin/events-delete/bootstrap

.PHONY: format
format:
	@go fmt ./internal/... ./pkg/... ./cmd/...

.PHONY: start
start:
	@echo "=> Starting service"
	@docker compose up -d
	@make format
	@make build
	@sls offline

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