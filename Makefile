NAME=weather
OUTPUT=cmd/$(NAME)

build:
	go build -o $(OUTPUT) cmd/main.go
test:
	go test ./...

install: build
	mv $(OUTPUT) $(GOPATH)/bin/

clean:
	rm -f $(OUTPUT)
