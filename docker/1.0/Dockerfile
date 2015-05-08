FROM buildpack-deps:jessie

RUN mkdir -p /var/lib/burlesque
WORKDIR /var/lib/burlesque

RUN apt-get update\
	&& apt-get install unzip libkyotocabinet16\
	&& rm -rf /var/lib/apt/lists/*\
	&& curl -SL https://github.com/KosyanMedia/burlesque/releases/download/v1.0.0/burlesque1.0.0.linux-amd64.zip -o burlesque.zip\
	&& unzip burlesque.zip\
	&& apt-get purge -y unzip\
	&& apt-get clean\
	&& rm burlesque.zip\
	&& chmod 755 burlesque\
	&& mv ./burlesque /usr/bin/

ENTRYPOINT /usr/bin/burlesque

EXPOSE 4401
