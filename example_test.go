package picam_test

import (
	"fmt"
	"log"

	"github.com/cgxeiji/picam"
)

func Example() {
	cam, err := picam.New(640, 480, picam.YUV)
	if err != nil {
		log.Fatal(err)
	}
	defer cam.Close()

	nFrames := 5
	fmt.Println("Reading", nFrames, "frames:")

	for i := 0; i < nFrames; i++ {
		// Get an image.Image
		img := cam.Read()

		/* do something with img */
		fmt.Println("got", img.Bounds().Size())

		// Or get a raw []uint8 slice
		raw := cam.ReadUint8()

		/* do something with img */
		fmt.Println("read", len(raw), "bytes")
	}

	// Output:
	// Reading 5 frames:
	// got (640,480)
	// read 460800 bytes
	// got (640,480)
	// read 460800 bytes
	// got (640,480)
	// read 460800 bytes
	// got (640,480)
	// read 460800 bytes
	// got (640,480)
	// read 460800 bytes
}
