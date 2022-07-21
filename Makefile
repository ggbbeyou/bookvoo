
test:
	go clean -testcache
	cd core/base && go test -v ./...
	cd user/assets && go test -v ./...
	# cd user/orders && go test -v ./...


depend:
	# https://juejin.cn/post/6844903609390333965	
	docker pull registry.docker-cn.com/swaggerapi/swagger-editor
	docker run --rm -p 9001:8080 swaggerapi/swagger-editor
	brew tap go-swagger/go-swagger
	brew install go-swagger