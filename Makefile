simple-bank-run:
	go run ./simple-bank-run/

simple-bank-createSQL:
	migrate create -ext sql -dir ./simple-bank/migrations/ -seq init_schema

simple-bank-up:
	migrate -path simple-bank/migrations -database "postgresql://postgres:12345@localhost:5432/simple_bank" -verbose up	

####============DOCKER=================####
postgres:
    docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345 -d postgres:16-

createdb:
    docker exec -it postgres16 createdb --username=postgres --owner=postgres simple_bank

dropdb:
    docker exec -it postgres16 dropdb simple_bank

.PHONY: simple-bank-run simple-bank-createSQL postgres createdb dropdb		