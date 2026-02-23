.PHONY: build install clean release-dry-run release

build:
	go build -o bl .

install:
	go install .

clean:
	rm -f bl

release-dry-run:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean
