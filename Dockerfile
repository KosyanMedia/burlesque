FROM buildpack-deps:jessie

ENV GOPATH /gocode
COPY . /gocode/src/github.com/KosyanMedia/burlesque
WORKDIR /gocode/src/github.com/KosyanMedia/burlesque

RUN apt-get update\
    && apt-get install -y golang libkyotocabinet16 libkyotocabinet16-dev\
    && rm -rf /var/lib/apt/lists/*\
    && go get -d -v\
  	&& go install -v\
  	&& go clean\
    && apt-get purge -y golang libkyotocabinet16-dev\
    && apt-get autoremove --purge -y\
    && apt-get clean\
  	&& ls /gocode/src/ | fgrep -v github.com | xargs rm -rf\
  	&& ls /gocode/src/github.com/ | fgrep -v KosyanMedia | xargs rm -rf

RUN mkdir /burlesque
WORKDIR /burlesque

ENTRYPOINT /gocode/bin/burlesque

EXPOSE 4401
