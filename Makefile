SWAG_ENTRY := internal/router.go
SWAG_OUT := docs

.PHONY: swagger test

swag:
	go run github.com/swaggo/swag/cmd/swag@latest init -g $(SWAG_ENTRY) -o $(SWAG_OUT)

test:
	go test ./... -count=1
