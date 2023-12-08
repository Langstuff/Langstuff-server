all: create_bin server langstuff

.PHONY: all create_bin server

create_bin:
	rm -rf bin
	mkdir bin

langstuff: cmd/langstuff.go
	go build -o ./bin/langstuff cmd/langstuff.go

cmd-%:
	go build -o ./bin/$* ./cmd/$*

server:
	go build -o ./bin/server ./server
