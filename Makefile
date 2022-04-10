include .env

launch-db:
	docker pull mysql:5.7.37
	docker run \
	--name mm14-mysql-database \
	-e MYSQL_ROOT_PASSWORD=${MYSQL_DB_PASSWORD} \
	-p ${MYSQL_DB_PORT}:3306 \
	-d mysql:5.7.37	

start-api: export INPUT_MODE=api

start-api:
	echo ${INPUT_MODE}
	go run main.go init_exchange kyiv_central
