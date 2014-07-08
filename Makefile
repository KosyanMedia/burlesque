.PHONY: all nuke

all:
	go get github.com/mattn/gom
	mkdir -p _vendor/bin
	gom install
	gofmt -w ./main.go
	gom build -o _vendor/bin/goqueue ./main.go

nuke:
	go clean -i