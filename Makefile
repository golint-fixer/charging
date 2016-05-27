all: test errcheck gofmt

test:
	go tool vet -test .

errcheck:
	errcheck ./...

gofmt:
	golint ./...
