.PHONY: prep
prep:
	docker run --name golang-todo-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres

.PHONY: migrate
migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

.PHONY: start
start:
	go run cmd/main.go