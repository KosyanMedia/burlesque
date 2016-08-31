FROM debian:jessie

ENV GOLANG_VERSION 1.6.2
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOPATH /
ENV GOROOT /go
ENV PATH $GOPATH/bin:/go/bin:$PATH

COPY . /src/github.com/KosyanMedia/burlesque
WORKDIR /src/github.com/KosyanMedia/burlesque

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates build-essential \
    libkyotocabinet16 libkyotocabinet16-dev curl git pkg-config git
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && tar -C / -xzf golang.tar.gz \
  && rm -rf golang.tar.gz \
  && go get -u github.com/kardianos/govendor \
  && govendor add +external && govendor get \
  && go install \
  && apt-get purge -y --auto-remove ca-certificates \
  && rm -rf /var/lib/apt/lists/* \
  && ln -s /bin/burlesque /bin/goqueue

CMD ["/bin/burlesque"]
EXPOSE 4401
