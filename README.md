# PiCam

[![Version](https://img.shields.io/github/v/tag/cgxeiji/picam?sort=semver)](https://github.com/cgxeiji/picam/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/cgxeiji/picam)](https://pkg.go.dev/github.com/cgxeiji/picam)
[![License](https://img.shields.io/github/license/cgxeiji/picam)](https://github.com/cgxeiji/picam/blob/master/LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/cgxeiji/picam)


PiCam is a Go wrapper to `raspiyuv` to get `[]uint8` and `image.Image` data of
the latests frame captured by the Raspberry Pi camera.

Under the hood, it executes:
```
$ raspiyuv --timeout 0 --timelapse 0
```
to get raw frames.

Currently, three image formats are available:
* picam.YUV
* picam.RGB
* picam.Gray

The time between frames measured on a Raspberry Pi Zero W is between `180ms` to
`210ms`.

## Why this library?

I wanted to avoid the dependency on GoCV to access the camera on a Raspberry
Pi to do real-time face detection.

## Example code

```go
package main

import (
	"image/png"
	"log"
	"os"

	"github.com/cgxeiji/picam"
)

func main() {
	cam, err := picam.New(640, 480, picam.YUV)
	if err != nil {
		log.Fatal(err)
	}
	defer cam.Close()

	img := cam.Read()

	f, err := os.Create("./image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Example real-time face detection

Using the [pigo](https://github.com/esimov/pigo) libray from esimov, it is
possible to do real-time face detection on a Raspberry Pi without depending on
OpenCV.

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cgxeiji/picam"
	pigo "github.com/esimov/pigo/core"
)

func main() {
	cam, err := picam.New(640, 480, picam.Gray)
	if err != nil {
		log.Fatal(err)
	}
	defer cam.Close()

	cParams := pigo.CascadeParams{
		MinSize:     90,
		MaxSize:     200,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Rows: cam.Height,
			Cols: cam.Width,
			Dim:  cam.Width,
		},
	}

	classifierFile, err := ioutil.ReadFile("./facefinder")
	if err != nil {
		log.Fatal(err)
	}

	p := pigo.NewPigo()
	classifier, err := p.Unpack(classifierFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting face detection")
	fmt.Println("Press Ctrl+C to stop")
	for {
		cParams.Pixels = cam.ReadUint8()
		faces := classifier.RunCascade(cParams, 0.0) // 0.0 is the angle
		faces = classifier.ClusterDetections(faces, 0.1)

		// Get the face with the highest confidence level
		var maxQ float32
		index := 0
		for i, face := range faces {
			if face.Q > maxQ {
				maxQ = face.Q
				index = i
			}
		}

		face := pigo.Detection{}
		if index < len(faces) {
			face = faces[index]
		}

		if face.Scale == 0 {
			// no face detected
			fmt.Printf("\rno face detected                                                 ")
			continue
		}

		x := face.Col - cam.Width/2
		y := -face.Row + cam.Height/2 // y is flipped

		fmt.Printf("\rface is (%d, %d) pixels from the center", x, y)
	}
}
```
