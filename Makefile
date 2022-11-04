

init:
	# admin 
	go get github.com/GoAdminGroup/adm
	go install github.com/GoAdminGroup/adm
	# wget -O admin/db/admin.sql https://gitee.com/go-admin/go-admin/raw/master/data/admin.sql 
	# wget -O admin/db/admin.pgsql https://gitee.com/go-admin/go-admin/raw/master/data/admin.pgsql 
	mysql -h127.0.0.1:13306 -uroot -proot -D test<admin/db/admin.sql

run-admin:
	cd admin/
	go run main.go

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
	cd clearing/ && go test -v ./...
	cd quotation/ && go test -v ./...
	go run main.go


docs:
	swag init

.PHONY: docs test debug run init