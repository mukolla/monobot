.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t monobot:v0.1 .

start-container:
	docker run --name talegram-monobot --env-file .env monobot:v0.1

push-dockerhub:
	docker build -t monobot:v0.1 .
	docker tag monobot:v0.1 mukolla/monobot:latest
	docker push mukolla/monobot:latest