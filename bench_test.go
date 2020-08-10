package picam

import (
	"fmt"
	"testing"
)

func BenchmarkPicam(b *testing.B) {
	benchmarks := []struct {
		format Format
	}{
		{YUV},
		{RGB},
		{Gray},
	}

	for _, bm := range benchmarks {
		b.Run(fmt.Sprintf("%s", bm.format), func(b *testing.B) {
			cam, err := New(640, 480, bm.format)
			if err != nil {
				b.Fatal(err)
			}
			defer cam.Close()
			_ = cam.ReadUint8()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = cam.ReadUint8()
			}
		})
	}
}
