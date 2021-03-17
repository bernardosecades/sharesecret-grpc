.PHONY: help dependencies up start stop restart status ps clean and execute tests

purge-secrets:
	docker-compose exec service bash -c "/bin/purge"
client-grpc-connection-example:
	docker-compose exec service bash -c "cd ./cmd/client && go build && ./client"
ps:
	docker-compose ps
up:
	docker-compose up --build
down:
	docker-compose down
test-coverage:
	docker-compose exec service bash -c "go clean -testcache  && go test  ./... -tags=unit,integration,e2e -coverprofile cover.out"
test-all:
	docker-compose exec service bash -c "go clean -testcache  && go test ./... -tags=unit,integration,e2e"
test-unit:
	docker-compose exec service bash -c "go clean -testcache  && go test ./... -tags=unit"
test-integration:
	docker-compose exec service bash -c "go clean -testcache  && go test ./... -tags=integration"
test-e2e:
	docker-compose exec service bash -c "go clean -testcache  && go test ./... -tags=e2e"
