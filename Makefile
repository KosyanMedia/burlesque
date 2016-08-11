.PHONY: github_release

docker_build:
	docker build -t aviasales/burlesque:latest .

docker_run:
	docker run --rm -p 4401:4401 aviasales/burlesque:latest

docker_tty:
	docker run --rm -p 4401:4401 -v `pwd`:/src/github.com/KosyanMedia/burlesque -ti aviasales/burlesque:latest /bin/bash

github_release:
	docker run --rm -v `pwd`:/src/github.com/KosyanMedia/burlesque aviasales/burlesque:latest /bin/bash -c "TAG=$(TAG) TOKEN=$(TOKEN) ./utils/github_release.sh"

build:
	go build --tags=leveldb main.go

test:
	docker run --rm -v `pwd`:/src/github.com/KosyanMedia/burlesque aviasales/burlesque:latest /bin/bash -c "cd clients/python/burlesque && python3 test.py"
