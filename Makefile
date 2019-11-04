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

