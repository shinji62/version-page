all: compile
compile: linux32 linux64 darwin64


linux32:
	GOARCH=386 GOOS=linux go build -o dist/linux/386/version-page_linux_386

linux64:
	GOARCH=amd64 GOOS=linux go build -o dist/linux/amd64/version-page_linux_amd64

darwin64:
	GOARCH=amd64 GOOS=darwin go build -o dist/darwin/amd64/version-page_darwin_amd64

clean:
	-rm -rf dist/*
	-rm -rf *.prof
