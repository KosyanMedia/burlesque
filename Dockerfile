FROM python:3.4.3

RUN curl -SL http://fallabs.com/kyotocabinet/pkg/kyotocabinet-1.2.76.tar.gz | tar -zxC /usr/src/\
    && cd /usr/src/kyotocabinet-1.2.76\
    && ./configure --prefix=/usr\
    && make\
    && make install\
    && cd /\
    && rm -rf /usr/src/kyotocabinet-1.2.76

ENV GOPATH /gocode
COPY . /gocode/src/github.com/KosyanMedia/burlesque
WORKDIR /gocode/src/github.com/KosyanMedia/burlesque

RUN apt-get update\
	&& apt-get install -y golang golang-doc golang-go golang-go-linux-amd64 golang-go.tools golang-src mercurial\
	&& apt-get clean\
	&& rm -rf /var/lib/apt/lists/*\
	&& go get -d -v\
	&& go install -v\
	&& go clean\
	&& ls /gocode/src/ | fgrep -v github.com | xargs rm -rf\
	&& ls /gocode/src/github.com/ | fgrep -v KosyanMedia | xargs rm -rf\
	&& apt-get purge -y golang golang-doc golang-go golang-go-linux-amd64 golang-go.tools golang-src mercurial

RUN mkdir /burlesque
WORKDIR /burlesque

ENTRYPOINT /gocode/bin/burlesque

EXPOSE 4401
