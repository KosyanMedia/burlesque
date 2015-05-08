FROM buildpack-deps:jessie

RUN curl -SL http://fallabs.com/kyotocabinet/pkg/kyotocabinet-1.2.76.tar.gz\
        | tar -zxC /usr/src/\
       && cd /usr/src/kyotocabinet-1.2.76\
       && ./configure --prefix=/usr\
       && make\
       && make install

RUN mkdir -p /var/lib/burlesque
WORKDIR /var/lib/burlesque

RUN apt-get update\
	&& apt-get install unzip\
	&& apt-get clean\
	&& rm -rf /var/lib/apt/lists/*\
	&& curl -SL https://github.com/KosyanMedia/burlesque/releases/download/v1.0.0/burlesque1.0.0.linux-amd64.zip -o burlesque.zip\
	&& unzip burlesque.zip\
	&& rm burlesque.zip\
	&& chmod 755 burlesque\
	&& mv ./burlesque /usr/bin/

ENTRYPOINT /usr/bin/burlesque

EXPOSE 4401
