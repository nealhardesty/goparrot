default: build

build:
	mkdir -p bin
	go build -o bin/goparrot ./...

clean:
	rm -f bin/goparrot
	rm -rf certs

run:
	go run . 

certs:
	mkdir -p certs
	openssl req -x509 -newkey rsa:4096 -nodes -keyout certs/key.pem -out certs/cert.pem -days 365 -subj "/C=US/CN=example.com"

runtls: certs
	go run . -cert certs/cert.pem -key certs/key.pem
