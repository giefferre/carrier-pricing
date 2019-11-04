# runs the whole project
.PHONY: start
start: bin generate-ssl-cert
	PROJECT_ROOT=`pwd` docker-compose up --build --force-recreate

# stops the docker compose
.PHONY: stop
stop:
	PROJECT_ROOT=`pwd` docker-compose stop

# tests code
.PHONY: tests
tests:
	go test -cover ./...

# builds the main binary
.PHONY: bin
bin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main cmd/main/main.go

# cleans the bin folder
.PHONY: clean
clean:
	cd bin && find . ! -name '.keep' -type f -exec rm -f {} +

# generates valid SSL certificates using github.com/FiloSottile/mkcert
generate-ssl-cert:
	go get github.com/FiloSottile/mkcert
	go run github.com/FiloSottile/mkcert --install -cert-file ./bin/localhost.cert -key-file ./bin/localhost.pem localhost 127.0.0.1 ::1
