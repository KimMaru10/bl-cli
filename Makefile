.PHONY: build install clean

build:
	go build -o bl .

install:
	go install .

clean:
	rm -f bl
