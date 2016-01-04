

default: test

test:
	go test -timeout 10s ./...

lint:
	gofmt -s -l timeglob
	go vet ./...

clean:
	go clean
