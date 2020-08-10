package picam

func gray(img []uint8) (g []uint8) {
	g = make([]uint8, len(img)/3)

	for i := 0; i < len(img)-3; i += 3 {
		g[i/3] = uint8(
			0.21*float64(img[i]) + // r
				0.72*float64(img[i+1]) + // g
				0.07*float64(img[i+2])) // b
	}

	return g
}
