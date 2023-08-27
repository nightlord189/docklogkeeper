buildg:
	docker build --no-cache -t "generator:latest" ./test/generator

rung:
	docker stop generator || true
	docker rm generator || true
	docker run --name generator generator:latest

swag:
	swag init --dir ./cmd/app --parseDependency --parseInternal

deploy:
	rm deploy.tar || true
	tar -cvf ./deploy.tar  ./*
	caprover deploy -t ./deploy.tar --host https://captain.app.tinygreencat.dev --caproverPassword ${CAPROVER_PASSWORD} --appName docklogkeeper
	rm deploy.tar