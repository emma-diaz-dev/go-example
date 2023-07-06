simple: 
	go run cmd/simple/main.go

docker-build-simple:
	docker build --no-cache -t cmd-simple -f ./Dockerfile.simple .

docker-run-simple:
	docker run --env-file .env -i cmd-simple

docker-deploy-simple: docker-build-simple docker-run-simple

docker-build-chan:
	docker build --no-cache -t cmd-chan -f ./Dockerfile.chan .

docker-run-chan:
	docker run --env-file .env -i cmd-chan

docker-deploy-chan: docker-build-chan docker-run-chan