run:  ##Run Code
	go run ./cmd/api/main.go


swag: ##Run Swagger
	swag init -g main.go -o ./cmd/api/docs

test:
	go test ./... -cover