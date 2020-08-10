package picam

import (
	"fmt"
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
				t.Errorf("got: %T, want: %T", got, ts.want)
			}
		})
	}
}
