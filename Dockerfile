FROM golang:1.8.3

RUN apt-get update \
  && apt-get install -y --no-install-recommends --fix-missing libleveldb-dev libleveldb1 libsnappy1

RUN go get -u github.com/kardianos/govendor
COPY . /go/src/github.com/KosyanMedia/burlesque
WORKDIR /go/src/github.com/KosyanMedia/burlesque
RUN make install

CMD ["/go/bin/burlesque"]
EXPOSE 4401
