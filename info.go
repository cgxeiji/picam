// Package picam is a Go wrapper to `raspiyuv` to get `[]uint8` and `image.Image` data of
// the latests frame captured by the Raspberry Pi camera.
//
// Under the hood, it executes:
//	$ raspiyuv --timeout 0 --timelapse 0
// to get raw frames.
//
// Currently, three image formats are available:
//	* picam.YUV
//	* picam.RGB
//	* picam.Gray
//
// The time between frames, measured on a Raspberry Pi Zero W, is between `180ms` to
// `210ms` for a `640x480` pixels image.
//
// If you want to test the speed in your system, run:
//	$ cd $(go env GOPATH)/src/github.com/cgxeiji/picam
//	$ go test -bench . -benchtime=10x
//
// This will take 10 frames and output the average time between each frame. Change
// `-benchtime=10x` to `100x` or `Nx` to change the number of frames to test.
package picam
