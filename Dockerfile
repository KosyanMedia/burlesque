FROM buildpack-deps:jessie

RUN apt-get update\
    && apt-get install -y build-essential zlib1g-dev pkg-config golang\
    && apt-get clean\
    && rm -rf /var/lib/apt/lists/*

RUN curl -SL http://fallabs.com/kyotocabinet/pkg/kyotocabinet-1.2.76.tar.gz\
        | tar -zxC /usr/src/\
       && cd /usr/src/kyotocabinet-1.2.76\
       && ./configure --enable-lzma --enable-static\
       && make\
       && make install

ENV PKG_CONFIG_PATH /usr/lib/pkgconfig:$PKG_CONFIG_PATH
ENV LD_LIBRARY_PATH /usr/local/lib:$LD_LIBRARY_PATH

ENV GOPATH /gocode
COPY . /gocode/src/github.com/KosyanMedia/burlesque

WORKDIR /gocode/src/github.com/KosyanMedia/burlesque

RUN go get -d -v
RUN go install -v
RUN go clean

RUN ls /gocode/src/ | fgrep -v github.com | xargs rm -rf
RUN ls /gocode/src/github.com/ | fgrep -v KosyanMedia | xargs rm -rf

RUN apt-get purge -y build-essential zlib1g-dev pkg-config golang \
	&& apt-get autoremove --purge -y\
	&& apt-get clean

RUN mkdir /burlesque
WORKDIR /burlesque

ENTRYPOINT /gocode/bin/burlesque

EXPOSE 4401
