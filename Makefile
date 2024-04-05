BINARY_NAME=myapp
CODE_DIR=./careerhub/apicomposer
CONTAINER_IMAGE_NAME=careerhub-api-composer


proto:
	@protoc mw/grpcmw/internal/*.proto  --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative  --go_opt=paths=source_relative  --proto_path=.

## test: runs all tests
test:	
	@echo "Testing..."
	@env POSTING_GRPC_ENDPOINT=${POSTING_GRPC_ENDPOINT} API_PORT=${API_PORT} SECRET_KEY=${SECRET_KEY} go test -p 1 -timeout 60s ./...
	

