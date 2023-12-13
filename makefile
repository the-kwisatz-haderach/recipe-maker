# https://github.com/jackc/tern
db-migrate: 
	~/go/bin/tern migrate \
	--config migrations/tern.conf \
	--migrations migrations

db-migrate-remote:
	~/go/bin/tern migrate \
	--conn-string "postgres://$(user):$(password)@recipe-maker.cctpacyerhjl.eu-north-1.rds.amazonaws.com:5432/$(database)?sslmode=require" \
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

build-image-proxy:
	docker build -t 044984945511.dkr.ecr.eu-north-1.amazonaws.com/recipe-maker-nginx:latest ./nginx

push-image:
	docker push 044984945511.dkr.ecr.eu-north-1.amazonaws.com/recipe-maker:latest

push-image-proxy:
	docker push 044984945511.dkr.ecr.eu-north-1.amazonaws.com/recipe-maker-nginx:latest
