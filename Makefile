include .env

launch-db:
	docker pull mysql:5.7.37
	docker run \
	--name mm14-mysql-database \
	-e MYSQL_ROOT_PASSWORD=${MYSQL_DB_PASSWORD} \
	-p ${MYSQL_DB_PORT}:3306 \
	-d mysql:5.7.37	

build:
	go build main.go

init: build
	@./main init ${SUPERUSER_PASSWORD}

test: build
	./main signup_user Cataleya Heaton cheaton11@mail.com 6666777788889999 password3 "2828-8282"

create-users:
	./main signup_user Ziggy Wilks zwilks@mail.com 1111111111111111 password1 ""
	./main signup_user Hawa Brookes hbrookes@mail.com 2222333344445555 password2 ""
	./main signup_user Cataleya Heaton cheaton@mail.com 6666777788889999 password3 ""

assign-broker-license:
	@./main assign_broker cheaton@mail.com 1919-9191

create-company:
	@./main signup_company PM "Prime Metal" hello@pm.com "+38 (066) 123-4567" password

add-commodity:
	@./main add_commodity zwilks@mail.com iron 200 $(shell ./main signin_company PM password)
	@./main add_commodity cheaton@mail.com silver 2.5 $(shell ./main signin_company PM password)
	@./main add_commodity cheaton@mail.com silver 2.5 $(shell ./main signin_company PM password)

check-commodity:
	@./main check_commodity $(shell ./main signin_user zwilks@mail.com password1)
	@./main check_commodity $(shell ./main signin_user hbrookes@mail.com password2)
	@./main check_commodity $(shell ./main signin_user cheaton@mail.com password3)

list-commodities:
	@./main list_commodities

add-order:
	@./main add_order sell iron 10 "" $(shell ./main signin_user zwilks@mail.com password1)
	@./main add_order sell iron 20 "" $(shell ./main signin_user zwilks@mail.com password1)
	@./main add_order sell iron 30 "" $(shell ./main signin_user zwilks@mail.com password1)
	@./main add_order buy iron 5 cheaton@mail.com $(shell ./main signin_user hbrookes@mail.com password2)

list-orders:
	@./main list_orders_my $(shell ./main signin_user zwilks@mail.com password1)
	@./main list_orders_my $(shell ./main signin_user hbrookes@mail.com password2)

list-orders-all:
	@./main list_orders_all $(shell ./main signin_user cheaton@mail.com password3)

cancel-order:
	@./main cancel_order 2 $(shell ./main signin_user zwilks@mail.com password1)

execute-order:
	@./main execute_order 2 4 1 $(shell ./main signin_user cheaton@mail.com password3)

cli: build 
	@./main cli

add-license: 
	@./main add_license "1122-3344" ${SUPERUSER_PASSWORD}

start: build init-db create-users assign-broker-license create-company add-commodity add-order cancel-order

