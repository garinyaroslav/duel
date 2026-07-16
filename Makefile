NAME = duel

all: linux windows

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME) .

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(NAME).exe .

clean:
	rm -f bin/$(NAME) bin/$(NAME).exe
