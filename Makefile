.PHONY: all clean

all: hostbin appbin nodebin

clean:
	rm -f -- nodebin appbin hostbin

nodebin:
	go build -o nodebin ./node

appbin:
	go build -o appbin ./app

hostbin:
	go build -o hostbin .
