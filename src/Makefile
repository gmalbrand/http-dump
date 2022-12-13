prepare:
	go mod download
	go get

build:
	go build -o http-dump main.go

run:
	HTTP_SERVER_PORT=4242 ./htt-dump

clean:
	rm -rf http-dump

all: prepare build run
