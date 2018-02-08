FROM golang:1.9.2

RUN apt-get update && apt-get install -y \
      wget \
      xz-utils

RUN wget http://johnvansickle.com/ffmpeg/releases/ffmpeg-release-64bit-static.tar.xz \
      && tar Jxvf ./ffmpeg-release-64bit-static.tar.xz \
      && cp ./ffmpeg*64bit-static/ffmpeg /usr/local/bin/

RUN go get github.com/soeyusuke/bananaman && \
      go get github.com/grafov/m3u8


CMD /bin/bash

