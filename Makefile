include .env
export

migrate-up:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATION} up

migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATION} down

generate:
	PSQL_DSN="$(PSQL_DSN)" \
	go run github.com/stephenafamo/bob/gen/bobgen-psql@latest \
	-c bobgen.yaml
