package picam_test

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cgxeiji/picam"
	pigo "github.com/esimov/pigo/core"
)

func Example_face_detection() {
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
