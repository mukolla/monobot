.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t monobot:v0.1 .

start-container:
	docker run --name talegram-monobot -p 8183:8183 --env-file .env monobot:v0.1