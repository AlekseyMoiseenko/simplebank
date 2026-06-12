DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
COMPOSE_FILE := docker-compose.yaml
MIGRATE_PATH = db/migration

migrateup:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/AlekseyMoiseenko/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

compose-up:
	docker compose -f $(COMPOSE_FILE) up -d

compose-down:
	docker compose -f $(COMPOSE_FILE) down

compose-ps:
	docker compose -f $(COMPOSE_FILE) ps

compose-logs:
	docker compose -f $(COMPOSE_FILE) logs -f --tail=200

.PHONY: migrateup migratedown migrateup1 migratedown1 db_docs db_schema sqlc test server mock proto compose-up compose-down compose-ps compose-logs