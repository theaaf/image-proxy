package svg

// #cgo LDFLAGS: -lrsvg-2
// #include <librsvg/rsvg.h>
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

type SVG struct {
	handle *C.RsvgHandle
}

type rsvgError struct {
	err *C.GError
}

func (e rsvgError) Error() string {
	return C.GoString(e.err.message)
}

func New(data []byte) (*SVG, error) {
	var err *C.GError
	svg := &SVG{
		handle: C.rsvg_handle_new_from_data((*C.uchar)(C.CBytes(data)), C.ulong(len(data)), &err),
	}
	if err != nil {
		return nil, &rsvgError{
			err: err,
		}
	}
	return svg, nil
}

func (svg *SVG) Rasterize(scalingFunction func(width, height int) (int, int)) (image.Image, error) {
	var dimensions C.RsvgDimensionData
	C.rsvg_handle_get_dimensions(svg.handle, &dimensions)

	width, height := int(dimensions.width), int(dimensions.height)
	if scalingFunction != nil {
		width, height = scalingFunction(int(dimensions.width), int(dimensions.height))
	}

	surface := C.cairo_image_surface_create(C.CAIRO_FORMAT_ARGB32, C.int(width), C.int(height))
	defer C.cairo_surface_destroy(surface)
	if status := C.cairo_surface_status(surface); status != C.CAIRO_STATUS_SUCCESS {
		return nil, fmt.Errorf("unable to create image surface: %v", status)
	}

	cairo := C.cairo_create(surface)
	if status := C.cairo_status(cairo); status != C.CAIRO_STATUS_SUCCESS {
		return nil, fmt.Errorf("unable to create cairo context: %v", status)
	}
	defer C.cairo_destroy(cairo)

	C.cairo_scale(cairo, C.double(width)/C.double(dimensions.width), C.double(height)/C.double(dimensions.height))

	if C.rsvg_handle_render_cairo(svg.handle, cairo) == 0 {
		return nil, fmt.Errorf("render error")
	}

	cData := C.cairo_image_surface_get_data(surface)
	if cData == nil {
		return nil, fmt.Errorf("unable to read image data")
	}

	stride := int(C.cairo_image_surface_get_stride(surface))
	data := C.GoBytes(unsafe.Pointer(cData), C.int(stride*height))
	result := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := data[y*stride+x*4 : y*stride+(x+1)*4]
			result.SetRGBA(x, y, color.RGBA{
				R: pixel[2],
				G: pixel[1],
				B: pixel[0],
				A: pixel[3],
			})
		}
	}

	return result, nil
}
