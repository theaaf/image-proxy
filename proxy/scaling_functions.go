package proxy

import (
	"github.com/disintegration/imaging"
	"regexp"
	"strings"
)

type CropType string

const (
	CropTypeCenter      CropType = "center"
	CropTypeLeft        CropType = "left"
	CropTypeRight       CropType = "right"
	CropTypeTopLeft     CropType = "top_left"
	CropTypeTop         CropType = "top"
	CropTypeTopRight    CropType = "top_right"
	CropTypeBottomLeft  CropType = "bottom_left"
	CropTypeBottom      CropType = "bottom"
	CropTypeBottomRight CropType = "bottom_right"
)

var cropTypeToAnchor = map[CropType]imaging.Anchor{
	CropTypeCenter:      imaging.Center,
	CropTypeLeft:        imaging.Left,
	CropTypeRight:       imaging.Right,
	CropTypeTopLeft:     imaging.TopLeft,
	CropTypeTopRight:    imaging.TopRight,
	CropTypeTop:         imaging.Top,
	CropTypeBottom:      imaging.Bottom,
	CropTypeBottomLeft:  imaging.BottomLeft,
	CropTypeBottomRight: imaging.BottomRight,
}

func (c *CropType) Anchor() *imaging.Anchor {
	if a, ok := cropTypeToAnchor[*c]; ok {
		return &a
	}
	return nil
}

var camelCaseRegex = regexp.MustCompile("(^|_).")

func (c *CropType) CamelCase() string {
	return camelCaseRegex.ReplaceAllStringFunc(string(*c), func(s string) string {
		return strings.TrimLeft(strings.ToUpper(s), "_")
	})
}

type ScalingFunction func(width, height int, allowUpscaling bool) (int, int)

func ScaleToFit(fitWidth, fitHeight int) ScalingFunction {
	return func(width, height int, allowUpscaling bool) (int, int) {
		if (width <= fitWidth && height <= fitHeight) && !allowUpscaling {
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
	return func(width, height int, allowUpscaling bool) (int, int) {
		if (width <= fillWidth || height <= fillHeight) && !allowUpscaling {
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
