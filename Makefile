include .env

launch-db:
	docker pull mysql:5.7.37
	docker run \
	--name mm14-mysql-database \
	-e MYSQL_ROOT_PASSWORD=${MYSQL_DB_PASSWORD} \
	-p ${MYSQL_DB_PORT}:3306 \
	-d mysql:5.7.37	

export INPUT_MODE=api
export JWT_USER_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJlbWFpbCI6InNiZW5kZXJAbWFpbC51YSIsImV4Y2giOiJLQ0UiLCJleHAiOjE2NTA3OTcxNjF9.fEdZqi-bt4C5xIrT70Z-dsYWmutS-YBbC2572n4rItY
export JWT_COMPANY_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTA3OTU5MDAsInRhZyI6IlBFUiIsInR5cGUiOiJjIn0.cipiQfjS3Pn-ybUR2UV9B04Noqt2bxqfQKyiiELT_D8
export JWT_SHIPCOMAPNY_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NTA3ODUzODksInRhZyI6Ik9NS0giLCJ0eXBlIjoicyJ9._g2AOyk-vrOcCgg3svFet2atE-NXcwMwwjypFkkzPko



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
	./main add_commodity KCE sbender@mail.ua iron 200 ${JWT_COMPANY_TOKEN}

check-commodity:
	./main check_commodity ${JWT_USER_TOKEN}

check-all-commodity:
	./main check_commodity_broker KCE ${JWT_USER_TOKEN}

init: init-market-db init-exchangers-db

signup: create-users create-company create-shipment-company

signin: 
	./main signin_user KCE sbender@mail.ua vashbatk0
	./main signin_company PER paprikaa
	./main signin_shipcompany OMKH omkhpass

list_commodities:
	./main list_commodities

add_order:
	./main add_order sell iron 500 sbender@mail.ua $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)
	./main add_order sell silver 10 "" $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)
	./main add_order buy cooper 500 "" $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)
	./main add_order buy iron 200 "" $(shell ./main signin_user KCE gkh@mail.ua topgetman123)
	./main add_order sell gold 500 sbender@mail.ua $(shell ./main signin_user KCE gkh@mail.ua topgetman123)
	./main add_order buy cooper 400 "" $(shell ./main signin_user KCE gkh@mail.ua topgetman123)

list_orders:
	./main list_orders true $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)
	./main list_orders true $(shell ./main signin_user KCE gkh@mail.ua topgetman123)

list_broker_orders:
	./main list_orders_all KCE $(shell ./main signin_user KCE sbender@mail.ua vashbatk0)

update_order:
	./main update_order 8 sell silver 100 "sbender@mail.ua" $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)

delete_order:
	./main delete_order 10 $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)

start: build delete_order

