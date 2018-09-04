FROM golang:latest

WORKDIR /go

RUN apt-get update && apt-get install -y \
      wget \
      xz-utils

RUN wget http://johnvansickle.com/ffmpeg/releases/ffmpeg-release-64bit-static.tar.xz \
      && tar Jxvf ./ffmpeg-release-64bit-static.tar.xz \
      && cp ./ffmpeg*64bit-static/ffmpeg /usr/local/bin/

RUN git clone https://github.com/soeyusuke/bananaman && \
    go get -u github.com/golang/dep/cmd/dep

RUN cd ./bananaman/bananamoon
RUN dep ensure

RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

RUN apt-get update -y --disableplugin=fastestmirror && \
    apt-get install -y epel-release --disableplugin=fastestmirror && \
    apt-get install -y --disableplugin=fastestmirror sudo cronie

RUN sed -i -e '/pam_loginuid.so/s/^/#/' /etc/pam.d/crond
ADD cron.d /etc/cron.d/
RUN chmod 0644 /etc/cron.d/*

CMD crond && tail -f /dev/null

CMD /bin/bash

