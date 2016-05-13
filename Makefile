docker_build:
	docker build -t aviasales/burlesque --force-rm --no-cache .

docker_run:
	docker run --rm -p 4001:4001 aviasales/burlesque

docker_tty:
	docker run --rm -p 4001:4001 -v `pwd`:/src/github.com/KosyanMedia/burlesque -ti aviasales/burlesque /bin/bash
