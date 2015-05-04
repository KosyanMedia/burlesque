FROM centos:6

RUN yum install -y tar unzip gcc gcc-c++ zlib-devel\
	&& curl -SL http://fallabs.com/kyotocabinet/pkg/kyotocabinet-1.2.76.tar.gz | tar -zxC /usr/src/\
	&& cd /usr/src/kyotocabinet-1.2.76\
	&& ./configure --prefix=/usr\
	&& make\
	&& make install\
	&& cd / && rm -rf /usr/src/kyotocabinet-1.2.76\
	&& curl -SL https://github.com/KosyanMedia/burlesque/releases/download/v1.0.0/burlesque1.0.0.linux-amd64.zip -o burlesque.zip\
	&& unzip burlesque.zip\
	&& rm burlesque.zip\
	&& chmod 755 burlesque\
	&& mv ./burlesque /usr/bin/\
	&& mkdir -p /var/lib/burlesque\
	&& yum -y history undo 2\
	&& yum clean all

ENV LD_LIBRARY_PATH=/usr/lib

WORKDIR /var/lib/burlesque

ENTRYPOINT /usr/bin/burlesque

EXPOSE 4401
