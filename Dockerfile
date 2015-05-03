FROM buildpack-deps:jessie

RUN apt-get update\
    && apt-get install -y build-essential zlib1g-dev pkg-config golang\
    && apt-get clean\
    && rm -rf /var/lib/apt/lists/*

RUN curl -SL http://fallabs.com/kyotocabinet/pkg/kyotocabinet-1.2.76.tar.gz\
        | tar -zxC /usr/src/\
       && cd /usr/src/kyotocabinet-1.2.76\
       && ./configure --prefix=/usr\
       && make\
       && make install

ENV PKG_CONFIG_PATH /usr/lib/pkgconfig:$PKG_CONFIG_PATH

ENV GOPATH /gocode
COPY . /gocode/src/github.com/KosyanMedia/burlesque

WORKDIR /gocode/src/github.com/KosyanMedia/burlesque

RUN go get -d -v\
	&& go install -v\
	&& go clean\
	&& ls /gocode/src/ | fgrep -v github.com | xargs rm -rf\
	&& ls /gocode/src/github.com/ | fgrep -v KosyanMedia | xargs rm -rf

RUN mkdir /burlesque
WORKDIR /burlesque

ENTRYPOINT /gocode/bin/burlesque

EXPOSE 4401
