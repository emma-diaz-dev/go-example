simple: 
	go run cmd/simple/main.go

docker-build:
	docker build --no-cache -t cmd-simple -f ./Dockerfile.simple .

docker-run:
	docker run -i cmd-simple

docker-deploy: docker-build docker-run