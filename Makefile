.PHONY: github_release

docker_build:
	docker build -t aviasales/burlesque:latest --force-rm --no-cache .

docker_run:
	docker run --rm -p 4001:4001 aviasales/burlesque:latest

docker_tty:
	docker run --rm -p 4001:4001 -v `pwd`:/src/github.com/KosyanMedia/burlesque -ti aviasales/burlesque:latest /bin/bash

github_release:
	docker run --rm -v `pwd`:/src/github.com/KosyanMedia/burlesque aviasales/burlesque:latest /bin/bash -c "TAG=$(TAG) TOKEN=$(TOKEN) ./utils/github_release.sh"
