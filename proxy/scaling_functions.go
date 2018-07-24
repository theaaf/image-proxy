package proxy

type ScalingFunction func(width, height int) (int, int)

func ScaleToFit(fitWidth, fitHeight int) ScalingFunction {
	return func(width, height int) (int, int) {
		if width <= fitWidth && height <= fitHeight {
			return width, height
		}
		xScale := float64(fitWidth) / float64(width)
		yScale := float64(fitHeight) / float64(height)
		if xScale < yScale {
			return fitWidth, int(float64(height) * xScale)
		}
		return int(float64(width) * yScale), fitHeight
	}
}

func ScaleToFill(fillWidth, fillHeight int) ScalingFunction {
	return func(width, height int) (int, int) {
		if width <= fillWidth || height <= fillHeight {
			return width, height
		}
		xScale := float64(fillWidth) / float64(width)
		yScale := float64(fillHeight) / float64(height)
		if xScale < yScale {
			return int(float64(width) * yScale), fillHeight
		}
		return fillWidth, int(float64(height) * xScale)
	}
}
