.PHONY: github_release

docker_build:
	docker build -t aviasales/burlesque .

docker_run:
	docker run --rm -p 4401:4401 aviasales/burlesque

docker_tty:
	docker run --rm -p 4401:4401 -v `pwd`:/go/src/github.com/KosyanMedia/burlesque -ti aviasales/burlesque /bin/bash

github_release:
	docker run --rm -v `pwd`:/go/src/github.com/KosyanMedia/burlesque aviasales/burlesque /bin/bash -c "make install && TAG=$(TAG) TOKEN=$(TOKEN) ./utils/github_release.sh"

install:
	go build --tags leveldb -ldflags "-s" -o /go/bin/burlesque main.go

test:
	docker run --rm -v `pwd`:/go/src/github.com/KosyanMedia/burlesque aviasales/burlesque /bin/bash -c "cd clients/python/burlesque && python3 test.py"
