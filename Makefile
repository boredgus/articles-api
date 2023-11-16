composeName:=user-articles-api
composeFile:=docker/docker-compose.yml

start:
	docker compose -p $(composeName) -f $(composeFile) --env-file .env up

restart:
	docker compose -p $(composeName) -f $(composeFile) stop 
	docker rm articles-database
	docker rm articles-api
	docker rmi $(composeName)-database
	docker rmi $(composeName)-server
	docker compose -p $(composeName) -f $(composeFile) --env-file .env up
