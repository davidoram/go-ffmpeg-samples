package main

import "github.com/3d0c/gmf"

func main() {

	url := "rtsp://81.25.188.121:554/axis-media/media.amp?videocodec=h264&compression=30&fps=10&resolution=4CIF&videokeyframeinterval=30&camera=1"

	ctx, err := gmf.NewInputCtx(url)
	if err != nil {
		panic(err)
	}
	println("ok")
	defer ctx.CloseInputAndRelease()
	//...

}
