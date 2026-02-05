.PHONY: build run tidy clean

BINARY := clscli
ifeq ($(OS),Windows_NT)
	BINARY := clscli.exe
endif

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY) $(ARGS)

tidy:
	go mod tidy

clean:
	rm -f $(BINARY)
