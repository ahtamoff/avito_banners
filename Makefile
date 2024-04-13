pg_container:
	docker run --name=banners-db -e POSTGRES_PASSWORD='postgres' -p 5436:5432 -d --rm postgres
init_database:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5436/banners-db?sslmode=disable' up