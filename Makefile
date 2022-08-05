
run:
	@swag init
	@go run main.go

test:
	go clean -testcache
	cd common && go test -v ./...
	cd base/symbols && go test -v ./...
	cd user/assets && go test -v ./...
	cd user/orders && go test -v ./...
	cd clearings/ && go test -v ./...


.PHONY: docs
docs:
	swag init