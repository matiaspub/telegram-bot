.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-container:
	docker build -t telegram-bot:0.1 .

start-container:
	docker rm telegram-bot; docker run --name telegram-bot -p 8080:80 telegram-bot:0.1