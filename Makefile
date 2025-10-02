all:
	go build src/main.go

dev:
	DEBUG=true go run src/main.go examples/04.js

