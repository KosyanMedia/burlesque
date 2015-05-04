FROM python:3.4.3 

RUN curl -SL http://fallabs.com/kyotocabinet/pkg/kyotocabinet-1.2.76.tar.gz\
        | tar -zxC /usr/src/\
       && cd /usr/src/kyotocabinet-1.2.76\
       && ./configure --prefix=/usr\
       && make\
       && make install\
       && cd / && rm -rf /usr/src/kyotocabinet-1.2.76

RUN apt-get update && apt-get install unzip && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN mkdir /burlesque/
WORKDIR /burlesque

RUN curl -SL https://github.com/KosyanMedia/burlesque/releases/download/v1.0.0/burlesque1.0.0.linux-amd64.zip -o burlesque.zip && unzip burlesque.zip

ENTRYPOINT /burlesque/burlesque

EXPOSE 4401
