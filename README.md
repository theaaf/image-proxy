# image-proxy [![Build Status](https://travis-ci.org/theaaf/image-proxy.svg?branch=master)](https://travis-ci.org/theaaf/image-proxy)

Image proxy that can rasterize and resize images on-the-fly.

https://aaf.engineering/the-alliance-image-proxy/

# API

The origin URL is specified by removing the protocol and appending it to the image proxy's host. For example:

`https://image-proxy.example.com/files.example.com/image.svg`

Options can be specified via query parameters:

* `rasterize` - If the source image is a vector image, the proxy will rasterize it.
* `fit=[WIDTH]x[HEIGHT]` - Scales down the source image to fit entirely within the given bounds.
* `fill=[WIDTH]x[HEIGHT]` - Scales down the source image to entirely fill the given bounds.
    * In addition to the `fill` parameter, `crop` can be specified to return an image of the exact width and height, cropped by the parameter. Options are:
       * `crop=center`
       * `crop=left`
       * `crop=right`
       * `crop=top_left`
       * `crop=top`
       * `crop=top_right`
       * `crop=bottom_left`
       * `crop=bottom`
       * `crop=bottom_right`
* `format=jpg` - Converts raster images to JPEG. Vector images may still be retreived unless `rasterize` is also specified.
    * `quality=[QUALITY]` can be used to control the quality of the JPEG encoding, where `[QUALITY]` is a number ranging from 1 to 100 (inclusive).

# Configuration

You can specify a whitelist for allowed hosts via environment variable:

```
IMAGE_PROXY_ALLOWED_HOSTS=foo.example.com,bar.example.com,*.baz.example.com
```

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

The easiest way to deploy is via AWS CloudFormation. If you have an AWS account, you can deploy the proxy by just clicking this button:

[![Launch](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)](https://console.aws.amazon.com/cloudformation/home#/stacks/new?stackName=image-proxy&templateURL=https://s3.amazonaws.com/aaf-platform-prod-image-proxy-packaging/template-packaged.yaml)

Alternatively, see the aws-sam directory for the do-it-yourself deployment.
