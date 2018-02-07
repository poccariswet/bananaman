# bananamoon

 This repository is able to get the file is timefree of radiko. <br> But if you are area is not *kanto*, you must set the premium member's mail and pass. <br>
 So, please set the environment variable.
 
 ```
 RADIKO_MAIL
 RADIKO_PASS
 ```
## Installation

``` sh
$ go get github.com/soeyusuke/bananaman/bananamoon
$ go get github.com/grafov/m3u8
```

 
## Requirements
- ffmpeg
- go
- [radiko](http://radiko.jp/)

## Usage
```
$ bananamoon -id="stationID" -s="radio start time" -file=filename 
```
<br>

ex)

```
$ bananamoon -id=TBS -s=20171111010000 -file=bananamoonGold
```

## Reference
- [Radikoのタイムフリーで番組名を指定して録音してみる](http://d.hatena.ne.jp/nyanonon/touch/20161211)
- [authToken](https://github.com/yyoshiki41/go-radiko/blob/master/auth.go)
- [radicast](https://github.com/soh335/radicast/blob/master/radiko.go)
- [m3u8](https://ja.wikipedia.org/wiki/M3U)
- [ffmpeg](https://hori-ryota.com/blog/ffmpeg-mp4-concatenate/)


