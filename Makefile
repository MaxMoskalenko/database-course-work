include .env

launch-db:
	docker pull mysql:5.7.37
	docker run \
	--name mm14-mysql-database \
	-e MYSQL_ROOT_PASSWORD=${MYSQL_DB_PASSWORD} \
	-p ${MYSQL_DB_PORT}:3306 \
	-d mysql:5.7.37	

export INPUT_MODE=api
export JWT_USER_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJlbWFpbCI6ImdraEBtYWlsLnVhIiwiZXhjaCI6IktDRSIsImV4cCI6MTY1MDIxMjgyM30.OxIuS4XFfrgH6jLLhogfrjcYbYr1FOAC2gQ9CAPAVtI
export JWT_COMPANY_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTAyMTMwNTcsInRhZyI6IlBFUiIsInR5cGUiOiJjIn0.88AszQcrIFLc3moiAhH7Pp_kTvLs41g0WUA7EaB1xCA
export JWT_SHIPCOMAPNY_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTAyMTI4MjMsInRhZyI6Ik9NS0giLCJ0eXBlIjoicyJ9.tGdhgh--KQtz6XGOWwaebJrlmLNUDGtXdZ6wsj6elwE

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

add-commodity:
	./main add_commodity KCE ocherry@mail.ua iron 200 ${JWT_COMPANY_TOKEN}

init: init-market-db init-exchangers-db

signup: create-users create-company create-shipment-company

signin: 
	./main signin_user KCE gkh@mail.ua topgetman123
	./main signin_company PER paprikaa
	./main signin_shipcompany OMKH omkhpass

# start: build init-market-db init-exchangers-db create-users
start: build init signup signin

