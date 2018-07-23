package proxy

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime"
	"net/http"
	"net/textproto"

	"github.com/nfnt/resize"

	"github.aaf.cloud/platform/imp/svg"
)

type Filter func(*Response) (*Response, *FilterError)

type FilterError struct {
	Error      error
	StatusCode int
}

func RasterizeFilter(scalingFunction ScalingFunction) Filter {
	return func(in *Response) (*Response, *FilterError) {
		contentType, _, _ := mime.ParseMediaType(in.Header.Get("Content-Type"))
		if contentType != "image/svg+xml" {
			return in, nil
		}

		data, err := ioutil.ReadAll(in.Body)
		if err != nil {
			return nil, &FilterError{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
		}

		svg, err := svg.New(data)
		if err != nil {
			return nil, &FilterError{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
		}

		image, err := svg.Rasterize(scalingFunction)
		if err != nil {
			return nil, &FilterError{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
		}

		buf := &bytes.Buffer{}
		if err := png.Encode(buf, image); err != nil {
			return nil, &FilterError{
				Error:      err,
				StatusCode: http.StatusInternalServerError,
			}
		}

		out := &Response{
			Header: make(http.Header),
			Body:   buf,
		}

		for k, v := range in.Header {
			out.Header[k] = v
		}
		out.Header.Set("Content-Type", "image/png")

		return out, nil
	}
}

func ContentTypeFilter(in *Response) (*Response, *FilterError) {
	contentType, _, _ := mime.ParseMediaType(in.Header.Get("Content-Type"))
	switch contentType {
	case "image/jpeg", "image/png", "image/svg+xml":
		return in, nil
	}

	return nil, &FilterError{
		Error:      fmt.Errorf("forbidden content-type"),
		StatusCode: http.StatusForbidden,
	}
}

func HeaderFilter(in *Response) (*Response, *FilterError) {
	out := &Response{
		Header: http.Header{},
		Body:   in.Body,
	}
	for _, header := range forwardedHeaders {
		canonical := textproto.CanonicalMIMEHeaderKey(header)
		out.Header[canonical] = in.Header[canonical]
	}
	out.Header.Set("Content-Security-Policy", "default-src 'none'; style-src 'self' 'unsafe-inline'")
	out.Header.Set("X-Content-Type-Options", "nosniff")
	return out, nil
}

func ScaleFilter(scalingFunction ScalingFunction) Filter {
	return func(in *Response) (*Response, *FilterError) {
		var img image.Image
		var err error

		contentType, _, _ := mime.ParseMediaType(in.Header.Get("Content-Type"))
		switch contentType {
		case "image/svg+xml":
			return in, nil
		case "image/jpeg":
			img, err = jpeg.Decode(in.Body)
		case "image/png":
			img, err = png.Decode(in.Body)
		default:
			return nil, &FilterError{
				Error:      fmt.Errorf("unsupported content-type"),
				StatusCode: http.StatusForbidden,
			}
		}

		if err != nil {
			return nil, &FilterError{
				Error:      fmt.Errorf("unable to decode image"),
				StatusCode: http.StatusBadRequest,
			}
		}

		width, height := scalingFunction(img.Bounds().Dx(), img.Bounds().Dy())
		if width != img.Bounds().Dx() || height != img.Bounds().Dy() {
			img = resize.Resize(uint(width), uint(height), img, resize.Bicubic)
		}

		buf := &bytes.Buffer{}

		switch contentType {
		case "image/jpeg":
			err = jpeg.Encode(buf, img, &jpeg.Options{
				Quality: 98,
			})
		case "image/png":
			err = png.Encode(buf, img)
		}

		if err != nil {
			return nil, &FilterError{
				Error:      fmt.Errorf("unable to encode image"),
				StatusCode: http.StatusBadRequest,
			}
		}

		return &Response{
			Header: in.Header,
			Body:   buf,
		}, nil
	}
}
