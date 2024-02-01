run:  ##Run Code
	go run ./cmd/api/main.go


swag: ##Run Swagger
	swag init -g main.go -o ./cmd/api/docs

test: ##test
	go test ./... -cover

##mockgenerationrepo && usecase
## mockgen -source pkg\repository\interface\user.go -destination pkg\mock\mockRepo\user_mock.go -package mockRepository
## mockgen -source pkg\usecase\interface\user.go -destination pkg\mock\mockUseCase\user_mock.go -package mockUseCase   