// What it does:
//
// Program captures video from a connected camera,
// then uses the CascadeClassifier to detect faces,
// compares faces to faces in whitelist,
// if the face is not whitelisted is compares the face with already stored faces and
// saves it if the same person's face is not already saved
// 
//
// How to run:
//
// snooper [camera ID] [classifier XML file]
//
// 		go run ./cmd/main.go 0 ./data/haarcascade_frontalface_default.xml
//

package main

import (
	"fmt"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tsnooper [camera ID] [classifier XML file]")
		return
	}

	// parse args
	deviceID := os.Args[1]
	xmlFile := os.Args[2]

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Snooper detector")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("Start reading device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
