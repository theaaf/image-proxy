package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Request struct {
	OriginURL *url.URL
	Filters   []Filter
}

func parseDimensions(s string) (*Dimensions, error) {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid dimensions")
	}
	width, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid width")
	}
	height, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid height")
	}
	return &Dimensions{
		Width:  int(width),
		Height: int(height),
	}, nil
}

// NewProxyRequestFromURL parses a URL for the origin url and additional options.
func NewRequestFromURL(url *url.URL) (*Request, error) {
	r := &Request{}

	var err error

	r.OriginURL, err = url.Parse("https://" + strings.TrimPrefix(url.Path, "/"))
	if err != nil {
		return nil, err
	}

	scalingOptions := &ScalingOptions{}

	if crop := url.Query().Get("crop"); crop != "" {
		cropType := CropType(crop)

		if _, ok := cropTypeToAnchor[cropType]; !ok {
			return nil, fmt.Errorf("invalid crop: %v", crop)
		} else {
			scalingOptions.Crop = &cropType
		}
	}

	var scalingFunction ScalingFunction
	if fit := url.Query().Get("fit"); fit != "" {
		if dim, err := parseDimensions(fit); err != nil {
			return nil, errors.Wrap(err, "invalid fit")
		} else {
			scalingFunction = ScaleToFit(dim.Width, dim.Height)
			scalingOptions.Fit = dim
		}
	} else if fill := url.Query().Get("fill"); fill != "" {
		if dim, err := parseDimensions(fill); err != nil {
			return nil, errors.Wrap(err, "invalid fill")
		} else {
			scalingFunction = ScaleToFill(dim.Width, dim.Height)
			scalingOptions.Fill = dim
		}
	}

	r.Filters = append(r.Filters, ContentTypeFilter)

	if _, ok := url.Query()["rasterize"]; ok {
		r.Filters = append(r.Filters, RasterizeFilter(scalingFunction))
	}

	if scalingOptions.IsValid() {
		r.Filters = append(r.Filters, ScalingFilter(scalingOptions))
	}

	r.Filters = append(r.Filters, HeaderFilter)

	return r, nil
}

type Response struct {
	Header http.Header
	Body   io.Reader
}

type Configuration struct {
	AllowedHosts []string
}

func (c *Configuration) LoadEnvironmentVariables() {
	if hosts := os.Getenv("IMAGE_PROXY_ALLOWED_HOSTS"); hosts != "" {
		c.AllowedHosts = strings.Split(hosts, ",")
		for i, host := range c.AllowedHosts {
			c.AllowedHosts[i] = strings.TrimSpace(host)
		}
	}
}

func (c *Configuration) AllowsHost(host string) bool {
	if len(c.AllowedHosts) == 0 {
		return true
	}
	for _, pattern := range c.AllowedHosts {
		if ok, err := filepath.Match(pattern, host); err == nil && ok {
			return true
		}
	}
	return false
}

func Proxy(config *Configuration, w http.ResponseWriter, r *Request) {
	if !config.AllowsHost(r.OriginURL.Host) {
		http.Error(w, "forbidden host", http.StatusForbidden)
		return
	}

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
