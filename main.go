package main

/*
 * To execute: go run github.com/davidoram/go-ffmpeg-samples/main.go
 */
import (
	"fmt"
	"os"

	"github.com/lazywei/go-opencv/opencv"
)

// Image Context is passed through ImageHandlerFunc processing chain
type ImgCtx struct {
	// Window for output
	win *opencv.Window
	// Position of the threshold button on the output window
	pos int

	// Frame counter
	framecnt int

	// Font
	font *opencv.Font

	// Text to display
	text string
}

// Interface for methods that manipulate images
type ImageHandlerFunc func(*ImgCtx, *opencv.IplImage) (*opencv.IplImage, error)

func NewImgCtx() *ImgCtx {
	ctx := new(ImgCtx)
	ctx.win = opencv.NewWindow("Go-OpenCV Webcam")

	ctx.pos = 1
	ctx.text = "Moo"
	ctx.win.CreateTrackbar("Thresh", 1, 100, func(pos int, param ...interface{}) {
		ctx.pos = pos
	})
	ctx.font = opencv.InitFont(opencv.CV_FONT_HERSHEY_DUPLEX, 1, 1, 0, 1, 8)

	return ctx
}

func (this *ImgCtx) Destroy() {
	this.win.Destroy()
}

func (imgCtx *ImgCtx) GrayScale(img *opencv.IplImage) (*opencv.IplImage, error) {

	w := img.Width()
	h := img.Height()

	// Create the output image
	cedge := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 3)
	// defer cedge.Release()

	// Convert to grayscale
	gray := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
	edge := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
	defer gray.Release()
	defer edge.Release()

	opencv.CvtColor(img, gray, opencv.CV_BGR2GRAY)

	opencv.Smooth(gray, edge, opencv.CV_BLUR, 3, 3, 0, 0)
	opencv.Not(gray, edge)

	// Run the edge detector on grayscale
	opencv.Canny(gray, edge, float64(imgCtx.pos), float64(imgCtx.pos*3), 3)

	opencv.Zero(cedge)
	// copy edge points
	opencv.Copy(img, cedge, edge)

	return cedge, nil
}

func (imgCtx *ImgCtx) AddText(img *opencv.IplImage) (*opencv.IplImage, error) {

	w := img.Width()
	h := img.Height()
	c := opencv.NewScalar(255, 255, 255, 0)
	pos := opencv.Point{w / 2, h / 2}
	imgCtx.font.PutText(img, fmt.Sprintf("Frame %d", imgCtx.framecnt), pos, c)
	return img, nil
}

func (imgCtx *ImgCtx) Display(img *opencv.IplImage) (*opencv.IplImage, error) {
	imgCtx.win.ShowImage(img)
	return img, nil
}

func main() {
	imgCtx := NewImgCtx()
	defer imgCtx.Destroy()

	cap := opencv.NewCameraCapture(0)
	if cap == nil {
		panic("can not open camera")
	}
	defer cap.Release()

	for {
		if cap.GrabFrame() {
			img := cap.RetrieveFrame(1)
			if img != nil {
				imgCtx.framecnt += 1
				// fmt.Println("Processing frame")
				//ProcessImage(img, imgCtx)
				// imgCtx.ProcessPipeline(pipeline, img)
				img1, err := imgCtx.GrayScale(img)
				if err != nil {
					fmt.Printf("error stage 1: %v\n", err)
					continue
				}
				defer img1.Release()

				img2, err := imgCtx.AddText(img1)
				if err != nil {
					fmt.Printf("error stage 2: %v\n", err)
					continue
				}
				defer img2.Release()

				img3, err := imgCtx.Display(img2)
				if err != nil {
					fmt.Printf("error stage 3: %v\n", err)
					continue
				}
				defer img3.Release()

			} else {
				fmt.Println("Image ins nil")
			}
		}

		// Press 'esc' to quit
		if key := opencv.WaitKey(10); key == 27 {
			os.Exit(0)
		}
	}
	opencv.WaitKey(0)
}

// func (ctx *ImgCtx) ProcessPipeline(pipeline []*ImageHandlerFunc, img *opencv.IplImage) error {
// 	images := make([]*opencv.IplImage, len(pipeline))
// 	ctx.framecnt += 1
// 	for i, handler := range pipeline {
// 		println("frame: %d, handler %d", ctx.framecnt, i)
// 		img, err := handler(ctx, img)
// 		if err != nil {
// 			println(err)
// 			return err
// 		}
// 		images[i] = img
// 		defer images[i].Release()
// 	}
// 	return nil
// }

// func ProcessImage(img *opencv.IplImage, ctx *ImgCtx) error {
// 	w := img.Width()
// 	h := img.Height()
//
// 	// Create the output image
// 	cedge := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 3)
// 	defer cedge.Release()
//
// 	// Convert to grayscale
// 	gray := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
// 	edge := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
// 	defer gray.Release()
// 	defer edge.Release()
//
// 	opencv.CvtColor(img, gray, opencv.CV_BGR2GRAY)
//
// 	opencv.Smooth(gray, edge, opencv.CV_BLUR, 3, 3, 0, 0)
// 	opencv.Not(gray, edge)
//
// 	// Run the edge detector on grayscale
// 	opencv.Canny(gray, edge, float64(ctx.pos), float64(ctx.pos*3), 3)
//
// 	opencv.Zero(cedge)
// 	// copy edge points
// 	opencv.Copy(img, cedge, edge)
//
// 	ctx.win.ShowImage(cedge)
// 	return nil
// }
