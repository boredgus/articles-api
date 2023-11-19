composeName:=user-articles-api
composeFile:=docker/docker-compose.yml

start:
	docker compose -p $(composeName) -f $(composeFile) --env-file .env up

clean: 
	docker compose -p $(composeName) -f $(composeFile) down
	docker rmi $(composeName)-database
	docker rmi $(composeName)-server

restart:
	make clean
	make start