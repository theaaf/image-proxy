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

	"github.aaf.cloud/platform/image-proxy/svg"
	"github.com/disintegration/imaging"
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
		header := http.Header{}
		for k, v := range in.Header {
			header[k] = v
		}
		// Make absolutely sure there's only one content-type header so nothing downstream
		// interprets the content type differently.
		header.Set("Content-Type", in.Header.Get("Content-Type"))
		return &Response{
			Header: header,
			Body:   in.Body,
		}, nil
	}

	return nil, &FilterError{
		Error:      fmt.Errorf("forbidden content-type"),
		StatusCode: http.StatusForbidden,
	}
}

var forwardedHeaders = []string{
	"Cache-Control",
	"Content-Language",
	"Content-Type",
	"Expires",
	"Last-Modified",
	"Pragma",
}

func HeaderFilter(in *Response) (*Response, *FilterError) {
	out := &Response{
		Header: http.Header{},
		Body:   in.Body,
	}
	for _, header := range forwardedHeaders {
		canonical := textproto.CanonicalMIMEHeaderKey(header)
		if v, ok := in.Header[canonical]; ok {
			out.Header[canonical] = v
		}
	}
	out.Header.Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	out.Header.Set("X-Content-Type-Options", "nosniff")
	return out, nil
}

type Dimensions struct {
	Width  int
	Height int
}

type ScalingOptions struct {
	Crop *CropType
	Fill *Dimensions
	Fit  *Dimensions
}

func (o *ScalingOptions) IsValid() bool {
	if o.Crop != nil {
		// o.Crop cannot be provided by itself
		return o.Fill != nil
	}

	return o.Fill != nil || o.Fit != nil
}

func ScalingFilter(opts *ScalingOptions) Filter {
	return func(in *Response) (*Response, *FilterError) {
		var img image.Image
		var err error
		var crop *imaging.Anchor

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

		switch {
		case opts.Fit != nil:
			img = imaging.Fit(img, opts.Fit.Width, opts.Fit.Height, imaging.CatmullRom)
		case opts.Fill != nil:
			if opts.Crop != nil {
				if anchor := opts.Crop.Anchor(); anchor != nil {
					crop = anchor
				}
			}

			if crop != nil {
				width := opts.Fill.Width
				height := opts.Fill.Height

				if img.Bounds().Dx() < width {
					width = img.Bounds().Dx()
				}
				if img.Bounds().Dy() < height {
					height = img.Bounds().Dy()
				}

				img = imaging.Fill(img, width, height, *crop, imaging.CatmullRom)
			} else {
				width, height := ScaleToFill(opts.Fill.Width, opts.Fill.Height)(img.Bounds().Dx(), img.Bounds().Dy(), false)
				img = imaging.Resize(img, width, height, imaging.CatmullRom)
			}
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

func JPEGFilter(quality int) Filter {
	return func(in *Response) (*Response, *FilterError) {
		var img image.Image
		var err error

		contentType, _, _ := mime.ParseMediaType(in.Header.Get("Content-Type"))
		switch contentType {
		case "image/svg+xml", "image/jpeg":
			return in, nil
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

		buf := &bytes.Buffer{}

		err = jpeg.Encode(buf, img, &jpeg.Options{
			Quality: quality,
		})

		if err != nil {
			return nil, &FilterError{
				Error:      fmt.Errorf("unable to encode image"),
				StatusCode: http.StatusBadRequest,
			}
		}

		out := &Response{
			Header: make(http.Header),
			Body:   buf,
		}

		for k, v := range in.Header {
			out.Header[k] = v
		}
		out.Header.Set("Content-Type", "image/jpeg")

		return out, nil
	}
}
