# https://github.com/jackc/tern
db-migrate: 
	~/go/bin/tern migrate \
	--config migrations/tern.conf \
	--migrations migrations

db-migrate-down: 
	~/go/bin/tern migrate --destination -1 \
	--config migrations/tern.conf \
	--migrations migrations

db-migrate-up: 
	~/go/bin/tern migrate --destination +1 \
	--config migrations/tern.conf \
	--migrations migrations

db-create-migration: 
	~/go/bin/tern new $(name) \
	--migrations migrations

# https://github.com/99designs/gqlgen
gql-generate:
	go run github.com/99designs/gqlgen generate

# Docker artifacts
build-image:
	docker build -t 044984945511.dkr.ecr.eu-north-1.amazonaws.com/recipe-maker:latest .

push-image:
	docker push 044984945511.dkr.ecr.eu-north-1.amazonaws.com/recipe-maker:latest
