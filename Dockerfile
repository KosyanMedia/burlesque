FROM golang:1.7.0

RUN apt-get update \
  && apt-get install -y --no-install-recommends --fix-missing libleveldb-dev libleveldb1 libsnappy1

COPY . /go/src/github.com/KosyanMedia/burlesque
WORKDIR /go/src/github.com/KosyanMedia/burlesque

RUN set -e \
  go get -u github.com/kardianos/govendor \
  make install 

CMD ["/go/bin/burlesque"]
EXPOSE 4401
