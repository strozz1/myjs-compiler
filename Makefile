all:
	go build src/main.go

dev: 
	DEBUG=true go run src/main.go $(FILE)

