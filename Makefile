# running
composeName:=user-articles-api
composeFile:=docker/docker-compose.yml
start:
	docker compose -p $(composeName) -f $(composeFile) --env-file .env up

clean: 
	docker compose -p $(composeName) -f $(composeFile) down
	docker rmi $(composeName)-database
	docker rmi $(composeName)-server

remove-data:
	docker volume rm $(composeName)_db

restart:
	docker stop articles-api
	docker rm articles-api
	docker rmi $(composeName)-server
	make start

# tests
generate-mocks:
	mockery --config=./config/.mockery.yaml

tests:
	go test ./... -v -coverprofile="coverage.txt" -covermode count
	go tool cover -func="coverage.txt"

show coverage:
	go tool cover -html="coverage.txt"

# docs
specFile:=docs/swagger.json
validate-docs:
	swagger validate $(specFile)

serve-docs:
	make validate-docs
	swagger serve --flavor=swagger --port=3033 $(specFile)