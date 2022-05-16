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

init-db:
	./main init

create-users:
	./main signup_user Ziggy Wilks zwilks@mail.com 1111111111111111 password1 ""
	./main signup_user Hawa Brookes hbrookes@mail.com 2222333344445555 password2 ""
	./main signup_user Cataleya Heaton cheaton@mail.com 6666777788889999 password3 1919-9191

assign-broker-license:
	./main assign_broker cheaton@mail.com 1919-9191

create-company:
	./main signup_company PM "Prime Metal" hello@pm.com "+38 (066) 123-4567" password

add-commodity:
	./main add_commodity zwilks@mail.com iron 200 $(shell ./main signin_company PM password)
	./main add_commodity cheaton@mail.com silver 2.5 $(shell ./main signin_company PM password)
	./main add_commodity cheaton@mail.com silver 2.5 $(shell ./main signin_company PM password)

check-commodity:
	./main check_commodity ${JWT_USER_TOKEN}

check-all-commodity:
	./main check_commodity_broker KCE ${JWT_USER_TOKEN}

signup: create-users create-company create-shipment-company

signin: 
	./main signin_user KCE sbender@mail.ua vashbatk0
	./main signin_company PER paprikaa
	./main signin_shipcompany NNP newnewnew

list_commodities:
	./main list_commodities

add_order:
	./main add_order sell iron 500 sbender@mail.ua $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)
	./main add_order sell silver 10 "" $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)
	./main add_order buy cooper 500 "" $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)
	./main add_order buy iron 200 "" $(shell ./main signin_user KCE gkh@mail.ua topgetman123)
	./main add_order sell gold 500 sbender@mail.ua $(shell ./main signin_user KCE gkh@mail.ua topgetman123)
	./main add_order buy cooper 400 "" $(shell ./main signin_user KCE gkh@mail.ua topgetman123)
	./main add_order buy iron 200 "" $(shell ./main signin_user LCE ashept@mail.ua iNg0dwetrust)
	./main add_order sell gold 500 mgrusha@mail.ua $(shell ./main signin_user LCE mgrusha@mail.ua 50hryvnas)
	./main add_order buy cooper 400 "" $(shell ./main signin_user LCE ashept@mail.ua iNg0dwetrust)

list_orders:
	./main list_orders true $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)
	./main list_orders true $(shell ./main signin_user KCE gkh@mail.ua topgetman123)

list_orders_native:
	./main list_orders_native KCE $(shell ./main signin_user KCE sbender@mail.ua vashbatk0)

list_orders_foreign:
	./main list_orders_foreign $(shell ./main signin_user KCE sbender@mail.ua vashbatk0)

update_order:
	./main update_order 8 sell silver 100 "sbender@mail.ua" $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)

delete_order:
	./main delete_order 10 $(shell ./main signin_user KCE ocherry@mail.ua jokewriter4)

create-race:
	./main create_race KCE LCE "2022-06-12 12:04" $(shell ./main signin_shipcompany NNP newnewnew)

read-races:
	./main read_races

update-race:
	./main update_race 1 LCE KCE "2023-06-12 12:04" $(shell ./main signin_shipcompany NNP newnewnew)

delete-race:
	./main delete_race 1 $(shell ./main signin_shipcompany NNP newnewnew)

finish-race:
	./main finish_race 1 $(shell ./main signin_shipcompany NNP newnewnew)

execute-native:
	./main execute_order 1 4 1 $(shell ./main signin_user KCE sbender@mail.ua vashbatk0)

execute-foreign:
	./main execute_foreign_order KCE 1 LCE 1 1 15 $(shell ./main signin_user KCE sbender@mail.ua vashbatk0)

init: init-db

start: build init

