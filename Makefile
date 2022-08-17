
run:
	@swag init
	@go run main.go

debug:
	@swag init
	@export CGO_CPPFLAGS="-Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wno-builtin-requires-header" 
	@CGO_ENABLED=1 go run -race -v main.go

test:
	go clean -testcache
	cd common && go test -v ./...
	cd base/symbols && go test -v ./...
	cd user/assets && go test -v ./...
	cd user/orders && go test -v ./...
	cd clearings/ && go test -v ./...
	cd  market/ && go test -v ./...


.PHONY: docs
docs:
	swag init