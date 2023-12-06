all: create_bin server cmd-extract_flashcards cmd-play_flashcards

.PHONY: all create_bin server

create_bin:
	rm -rf bin
	mkdir bin

cmd-%:
	go build -o ./bin/$* ./cmd/$*

server:
	go build -o ./bin/server ./server
