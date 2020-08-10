package picam

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"math"
	"os/exec"
	"strconv"
	"sync"
)

// Camera is a struct that stores camera information.
type Camera struct {
	cmd *exec.Cmd
	// Width sets the width of the image
	// Height sets the height of the image
	Width, Height int
	frame         <-chan []uint8
	format        Format
	done          chan struct{}
	ws            *sync.WaitGroup
}

// Format is the type of image that picam will output.
type Format uint8

//go:generate stringer -type=Format
const (
	// YUV 420 color format.
	YUV Format = iota
	// RGB color format.
	RGB
	// Gray color format.
	Gray
)

// New initializes and starts a raspiyuv process to capture RGB frames.
func New(width, height int, format Format) (*Camera, error) {
	args := []string{
		"--burst",
		"--width", strconv.Itoa(width),
		"--height", strconv.Itoa(height),
		"--timeout", "0",
		"--timelapse", "0",
		"--nopreview",
	}

	var img []uint8
	switch format {
	case RGB:
		args = append(args, "--rgb")
		img = make([]uint8, width*height*3)
	case Gray:
		args = append(args, "--luma")
		img = make([]uint8, width*height)
	default:
		w, h := roundUp(width, 32), roundUp(height, 16)
		img = make([]uint8, w*h+w*h/2)
	}

	args = append(args, []string{"--output", "-"}...)

	cmd := exec.Command("raspiyuv", args...)

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("picam: cannot create out pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("picam: unable to start picam: %w", err)
	}

	frame := make(chan []uint8)
	done := make(chan struct{})
	var ws sync.WaitGroup
	ws.Add(1)
	go func() {
		defer ws.Done()
		defer close(frame)
		defer cmd.Process.Kill()
		r := bufio.NewReader(out)
		for {
			_, _ = io.ReadFull(r, img)
			select {
			case <-done:
				return
			case frame <- img:
			default:
			}
		}
	}()

	return &Camera{
		Width:  width,
		Height: height,
		cmd:    cmd,
		frame:  frame,
		format: format,
		done:   done,
		ws:     &ws,
	}, nil
}

// Close closes picam.
func (c *Camera) Close() {
	close(c.done)
	c.ws.Wait()
}

// Read returns an image.Image interface of the last frame.
//
//	cam, _ := picam.New(width, height, format)
//	img := cam.Read()
//
// The type returned depends on the format passed at picam.New():
//
//	format        type(img)
//	----------    ---------------
//	picam.YUV  -> image.YCbCr 420
//	picam.RGB  -> image.NRGBA
//	picam.Gray -> image.Gray
func (c *Camera) Read() (img image.Image) {
	size := image.Rect(0, 0, c.Width, c.Height)
	switch c.format {
	case RGB:
		rgba := image.NewNRGBA(size)
		pixels := make([]uint8, c.Width*c.Height*4)
		rgb := <-c.frame
		for i, idx := 0, 0; i < len(rgb); i++ {
			pixels[idx] = rgb[i]
			idx++
			if i%3 == 2 {
				pixels[idx] = 255
				idx++
			}
		}
		rgba.Pix = pixels
		img = rgba
	case Gray:
		gray := image.NewGray(size)
		gray.Pix = <-c.frame
		img = gray
	default:
		yuv := image.NewYCbCr(size, image.YCbCrSubsampleRatio420)

		yRange := roundUp(c.Width, 32) * roundUp(c.Height, 16)
		uvRange := yRange / 4

		frame := <-c.frame
		yuv.Y = frame[0:yRange]
		yuv.Cb = frame[yRange : uvRange+yRange]
		yuv.Cr = frame[yRange+uvRange : uvRange*2+yRange]
		img = yuv
	}

	return img
}

func roundUp(value, multiple int) int {
	return int(math.Ceil(float64(value)/float64(multiple))) * multiple
}

// ReadUint8 returns the raw uint8 values of the last frame.
//
//	cam, _ := picam.New(width, height, format)
//	raw := cam.ReadUint8()
//
// The size of the slice returned depends on the format and dimensions passed at picam.New():
//
//	format        len(raw)
//	----------    ----------------------------------------------------------
//	picam.YUV  -> roundUpMultiple32(width) * roundUpMultiple16(height) * 1.5
//	picam.RGB  -> (width * height) * 3
//	picam.Gray -> width * height
func (c *Camera) ReadUint8() (img []uint8) {
	return <-c.frame
}
