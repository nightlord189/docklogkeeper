buildg:
	docker build --no-cache -t "bibgen:latest" ./test/bibgen

rung:
	docker stop bibgen || true
	docker rm bibgen || true
	docker run --name bibgen bibgen:latest

buildg2:
	docker build --no-cache -t "randomgen:latest" ./test/randomgen

rung2:
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