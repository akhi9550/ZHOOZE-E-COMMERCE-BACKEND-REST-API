run:  ##Run Code
	go run ./cmd/api/main.go


swag: ##Run Swagger
	swag init -g cmd/api/main.go -o ./cmd/api/docs


test: ##test
	go test ./... -cover

test-coverage: ## Run tests and generate coverage file
	go test ./... -coverprofile=code-coverage.out
	go tool cover -html=code-coverage.out

wire: ## Generate wire_gen.go
	cd pkg/di && wire

mock: ##make mock files using mockgen
	mockgen -source pkg\repository\interface\user.go -destination pkg\repository\mock\user_mock.go -package mock
	mockgen -source pkg\repository\interface\order.go -destination pkg\repository\mock\order_mock.go -package mock
	mockgen -source pkg\usecase\interface\user.go -destination pkg\usecase\mock\user_mock.go -package mock

