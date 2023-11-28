run:
	go run ./cmd/api/main.go


swag: 
	swag init -g main.go -o ./cmd/docs