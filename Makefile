
postgres:
	docker run --name bestwallet_db_1 -p 5434:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres

createdb:
	docker exec -it bestwallet_db_1 createdb --username=postgres bestwallet_db

dropdb:
	docker exec -it bestwallet_db_1 dropdb --username=postgres bestwallet_db

start:
	docker start bestwallet_db_1

stop:
	docker stop bestwallet_db_1

.PHONY: bestwallet_db_1 createdb dropdb
