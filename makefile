# https://github.com/jackc/tern
db-migrate: 
	~/go/bin/tern migrate \
	--config migrations/tern.conf \
	--migrations migrations

db-migrate-rerun-latest: 
	~/go/bin/tern migrate --destination -+1 \
	--config migrations/tern.conf \
	--migrations migrations

db-create-migration: 
	~/go/bin/tern new $(name) \
	--config migrations/tern.conf \
	--migrations migrations

# https://github.com/99designs/gqlgen
gql-generate:
	go run github.com/99designs/gqlgen generate