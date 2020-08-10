package picam

import (
	"fmt"
	"image"
	"image/color"
	"testing"
)

func TestNew(t *testing.T) {
	cam, err := New(640, 480, YUV)
	if err != nil {
		t.Fatal(err)
	}
	defer cam.Close()

}

func TestRead(t *testing.T) {
	c := color.RGBA{}
	tests := []struct {
		format Format
		want   color.Color
	}{
		{YUV, color.YCbCrModel.Convert(c)},
		{RGB, color.NRGBAModel.Convert(c)},
		{Gray, color.GrayModel.Convert(c)},
	}

	for _, ts := range tests {
		t.Run(fmt.Sprintf("%s", ts.format), func(t *testing.T) {
			cam, err := New(640, 480, ts.format)
			if err != nil {
				t.Fatal(err)
			}
			defer cam.Close()

			img := cam.Read()

			got := img.ColorModel().Convert(c)
			if got != ts.want {
				t.Errorf("got: %T, want: %T", got, ts.want)
			}
		})
	}
}

func TestReadUint8(t *testing.T) {
	w, h := 640, 480
	tests := []struct {
		format Format
		want   int // byte size
	}{
		{YUV, w*h + w*h/2},
		{RGB, w * h * 3},
		{Gray, w * h},
	}

	for _, ts := range tests {
		t.Run(fmt.Sprintf("%s", ts.format), func(t *testing.T) {
			cam, err := New(640, 480, ts.format)
			if err != nil {
				t.Fatal(err)
			}
			defer cam.Close()

			img := cam.ReadUint8()

			got := len(img)
			if got != ts.want {
				t.Errorf("got: %d, want: %d", got, ts.want)
			}
		})
	}
}

func TestReadUint8_Sizes(t *testing.T) {
	tests := []struct {
		format        Format
		width, height int
		want          int // byte size
	}{
		{
			YUV,
			320, 240,
			320*240 + 320*240/2,
		},
		{
			YUV,
			100, 100,
			128*112 + 128*112/2,
		},
		{
			RGB,
			320, 240,
			320 * 240 * 3,
		},
		{
			Gray,
			320, 240,
			320 * 240,
		},
	}

	for _, ts := range tests {
		t.Run(fmt.Sprintf("%s (%d,%d)", ts.format, ts.width, ts.height), func(t *testing.T) {
			cam, err := New(ts.width, ts.height, ts.format)
			if err != nil {
				t.Fatal(err)
			}
			defer cam.Close()

			raw := cam.ReadUint8()

			got := len(raw)
			if got != ts.want {
				t.Errorf("got: %d, want: %d", got, ts.want)
			}

			img := cam.Read()
			gotP := img.Bounds().Size()
			wantP := image.Point{ts.width, ts.height}
			if gotP != wantP {
				t.Errorf("got: %+v, want: %+v", gotP, wantP)
			}
		})
	}
}
