# OpenCV experiments

## Setup

To pick up my forked changes of `github.com/lazywei/go-opencv` (until PR merged):
```
$ go get github.com/lazywei/go-opencv/opencv
$ cd src/github.com/lazywei/go-opencv/
$ git remote add fork https://github.com/davidoram/go-opencv.git
$ git fetch fork
$ git checkout fork/master
```

To test that fork:
```
cd src/github.com/lazywei/go-opencv/
go test opencv/cxcore_test.go
```

To run the program
`go run src/github.com/davidoram/go-ffmpeg-samples/main.go`






# FFMPEG experiments

Setup
`brew install ffmpeg

brew reinstall ffmpeg --with-ffplay --with-tools --with-openssl --with-rtmpdump --with-libssh --with-libass --with-freetype --with-fontconfig
`

To create a stream & show it:

```
Open with vlc 'udp://@:1234'

# Raw
ffmpeg -f avfoundation -i "FaceTime" -vcodec libx264 -tune zerolatency \
-b 900k -f mpegts udp://localhost:4242

# With video filter
ffmpeg -f avfoundation -i "FaceTime"  -vf "vflip" -vcodec libx264 -tune zerolatency -b 200k -f mpegts udp://127.0.0.1:1234

ffmpeg -f avfoundation -i "FaceTime" -vf "vflip"  -vcodec libx264 -tune zerolatency -b 500k -f mpegts udp://127.0.0.1:1234

# Drop size down to 640 x 480
ffmpeg -f avfoundation -i "FaceTime" -vf "vflip" -s 640x480 -vcodec libx264 -tune zerolatency -b 500k -f mpegts udp://127.0.0.1:1234
```

Design idea:
- Implement a generic plugin platform design that can extend ffmpeg
  - Eg: Passes images via IPC (eg: http://zeromq.org/) to another process
        Every Nth frame
        Next Frame when asked
  - Separate process does image processing & passes them back to be interspersed in the stream


To Build
```
go run src/github.com/davidoram/go-ffmpeg-samples/main.go -input=src/in.mov
```

Testing changes
go test github.com/davidoram/gmf

https://trac.ffmpeg.org/wiki/FilteringGuide
http://wiki.multimedia.cx/index.php?title=FFmpeg_filter_howto
http://comments.gmane.org/gmane.comp.video.ffmpeg.devel/149999 - Not likely to implement plugin/filters
http://git.videolan.org/?p=ffmpeg.git;a=commitdiff;h=3250231a0292d716afd9d1ad25fc39bacda17f67 - commit of a new filter


Solution is to:
- Fork ffmpeg
- Write glue code that allows filters to be built in golang
-

