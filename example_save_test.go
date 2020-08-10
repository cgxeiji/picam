package picam_test

import (
	"image/png"
	"log"
	"os"

	"github.com/cgxeiji/picam"
)

func Example_save() {
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
