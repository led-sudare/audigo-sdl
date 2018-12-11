package util

func LinearTransF(x float64, sMin float64, sMax float64, dMin float64, dMax float64) float64 {
	// validate
	if sMin == sMax {
		return (dMax - dMin) / 2
	}
	if dMin == dMax {
		return dMin
	}

	portion := (x - sMin) * (dMax - dMin) / (sMax - sMin)
	res := portion + dMin

	if res < dMin {
		return dMin
	}
	if res > dMax {
		return dMax
	}
	return res
}
