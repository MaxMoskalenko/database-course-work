include .env

launch-db:
	docker pull mysql:5.7.37
	docker run \
	--name mm14-mysql-database \
	-e MYSQL_ROOT_PASSWORD=${MYSQL_DB_PASSWORD} \
	-p ${MYSQL_DB_PORT}:3306 \
	-d mysql:5.7.37	

export INPUT_MODE=api

build:
	go build main.go

init-market-db:
	./main init

init-exchangers-db:
	./main init_exchange kyiv_central_ex "Kyiv Commodity Exchange" KCE

create-users:
	./main signup_user KCE Elon Tusk itusk@yahoo.eu password false
	./main signup_broker KCE Elon Lux ilux@yahoo.eu password true "01234-itsli-cense-56789"

get-user: export JWT_TOKEN=`./main signin_user KCE itusk@yahoo.eu password`

get-user:
	./main test ${JWT_TOKEN}

# start: build init-market-db init-exchangers-db create-users
start: build get-user

