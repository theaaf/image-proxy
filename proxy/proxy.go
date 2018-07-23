package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Request struct {
	OriginURL *url.URL
	Filters   []Filter
}

func parseDimensions(s string) (int, int, error) {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid dimensions")
	}
	width, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width")
	}
	height, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height")
	}
	return int(width), int(height), nil
}

// NewProxyRequestFromURL parses a URL for the origin url and additional options.
func NewRequestFromURL(url *url.URL) (*Request, error) {
	r := &Request{}

	var err error

	r.OriginURL, err = url.Parse("https://" + strings.TrimPrefix(url.Path, "/"))
	if err != nil {
		return nil, err
	}

	var scalingFunction ScalingFunction
	if fit := url.Query().Get("fit"); fit != "" {
		if width, height, err := parseDimensions(fit); err != nil {
			return nil, errors.Wrap(err, "invalid fit")
		} else {
			scalingFunction = ScaleToFit(width, height)
		}
	} else if fill := url.Query().Get("fill"); fill != "" {
		if width, height, err := parseDimensions(fill); err != nil {
			return nil, errors.Wrap(err, "invalid fill")
		} else {
			scalingFunction = ScaleToFill(width, height)
		}
	}

	r.Filters = append(r.Filters, ContentTypeFilter)

	if _, ok := url.Query()["rasterize"]; ok {
		r.Filters = append(r.Filters, RasterizeFilter(scalingFunction))
	}

	if scalingFunction != nil {
		r.Filters = append(r.Filters, ScaleFilter(scalingFunction))
	}

	r.Filters = append(r.Filters, HeaderFilter)

	return r, nil
}

var forwardedHeaders = []string{
	"Cache-Control",
	"Content-Language",
	"Content-Type",
	"Expires",
	"Last-Modified",
	"Pragma",
}

type Response struct {
	Header http.Header
	Body   io.Reader
}

func Proxy(w http.ResponseWriter, r *Request) {
	resp, err := http.Get(r.OriginURL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode >= 400 {
			http.Error(w, fmt.Sprintf("upstream error: %v", resp.StatusCode), resp.StatusCode)
		} else {
			http.Error(w, fmt.Sprintf("upstream status code: %v", resp.StatusCode), http.StatusBadGateway)
		}
		return
	}

	out := &Response{
		Header: resp.Header,
		Body:   resp.Body,
	}

	for _, filter := range r.Filters {
		var err *FilterError
		out, err = filter(out)
		if err != nil {
			http.Error(w, err.Error.Error(), err.StatusCode)
			return
		}
	}

	for k, v := range out.Header {
		w.Header()[k] = v
	}
	io.Copy(w, out.Body)
}
