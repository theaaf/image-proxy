# image-proxy

Image proxy that can rasterize and resize images on-the-fly.

# API

The origin URL is specified by removing the protocol and appending it to the image proxy's host. For example:

`https://image-proxy.example.com/files.example.com/image.svg`

Options can be specified via query parameters:

* `rasterize` - If the source image is a vector image, the proxy will rasterize it.
* `fit=[WIDTH]x[HEIGHT]` - Scales down the source image to fit entirely within the given bounds.
* `fill=[WIDTH]x[HEIGHT]` - Scales down the source image to entirely fill the given bounds.

# Development

* Install librsvg. On macOS, you can `brew install librsvg`.
* Add librsvg include and lib paths to `CGO_CFLAGS` and `CGO_LDFLAGS` if necessary. For example, on macOS:

  ```
  export PKG_CONFIG_PATH=$(brew --prefix librsvg)
  export CGO_CFLAGS=$(pkg-config --cflags librsvg-2.0)
  export CGO_LDFLAGS=$(pkg-config --libs librsvg-2.0)
  ```
* Test and run via standard Go commands: `dep ensure`, `go test ./...`, `go run main.go`, etc.

# Deployment

The easiest way to deploy is via the AWS Serverless Application Model (SAM). See the aws-sam directory for details.
