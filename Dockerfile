FROM ubuntu:latest

ENV GOLANG_VERSION 1.5.3
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 43afe0c5017e502630b1aea4d44b8a7f059bf60d7f29dfd58db454d4e4e0ae53
ENV GO15VENDOREXPERIMENT 1

COPY . /src/github.com/KosyanMedia/burlesque
WORKDIR /src/github.com/KosyanMedia/burlesque

RUN apt-get update \
  && apt-get install -y libkyotocabinet16 libkyotocabinet16-dev curl git pkg-config \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get -v -d && go build -race -o /burlesque
ENTRYPOINT  "/burlesque"
EXPOSE 4401
