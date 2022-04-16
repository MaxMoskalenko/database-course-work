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
	./main signup_user KCE Goddamn Khmelnytskyi gkh@mail.ua topgetman123 false
	./main signup_user KCE Ostap Cherry ocherry@mail.ua jokewriter4 false
	./main signup_broker KCE Stephan Bender sbender@mail.ua vashbatk0 true "01234-itsli-cense-56789"

create-company:
	./main signup_company PER "Paprika Journal" paprikaa
	./main signup_company GET "TOV Getmanchyna" topgetman

create-shipment-company:
	./main signup_shipcompany OMKH "Odesa-Mykolaiv-Kherson Monopoly" omkhpass

get-user: export JWT_TOKEN=`./main signin_user KCE itusk@yahoo.eu password`

get-user:
	./main test ${JWT_TOKEN}

test-signup: create-users create-company create-shipment-company

test-signin: 
	./main signin_user KCE gkh@mail.ua topgetman123
	./main signin_company PER paprikaa
	./main signin_shipcompany OMKH omkhpass

# start: build init-market-db init-exchangers-db create-users
start: build test-signin

