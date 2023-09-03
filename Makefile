run:
	docker build --no-cache -t "nightlord189/docklogkeeper:latest" .
	docker stop docklogkeeper || true
	docker rm docklogkeeper || true
	docker run --name docklogkeeper -d -v /var/run/docker.sock:/var/run/docker.sock -v docklogkeeper:/logs -p 3010:3010 nightlord189/docklogkeeper:latest


.PHONY: build
publish:
	docker buildx build --no-cache --push \
      --platform linux/arm64/v8,linux/amd64 \
      --tag $(image) \
      .

runbib:
	docker build --no-cache -t "bibgen:latest" ./test/bibgen
	docker stop bibgen || true
	docker rm bibgen || true
	docker run --name bibgen bibgen:latest

runrandom:
	docker build --no-cache -t "randomgen:latest" ./test/randomgen
	docker stop randomgen || true
	docker rm randomgen || true
	docker run --name randomgen -d randomgen:latest

swag:
	swag init --dir ./cmd/app --parseDependency --parseInternal

deploy:
	rm deploy.tar || true
	tar -cvf ./deploy.tar  ./*
	caprover deploy -t ./deploy.tar --host https://captain.app.tinygreencat.dev --caproverPassword ${CAPROVER_PASSWORD} --appName docklogkeeper
	rm deploy.tar

migrate-new:
	goose -s -dir configs/migrations/local create $(name) sql

migrate:
	goose -dir configs/migrations/local sqlite3 ./logs.db up

migrate-down:
	goose -dir configs/migrations/local sqlite3 ./logs.db down

migrate-reset:
	goose -dir configs/migrations/local sqlite3 ./logs.db reset