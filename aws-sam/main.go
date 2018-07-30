package main

import (
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.aaf.cloud/platform/image-proxy/proxy"
)

type Request struct {
	Path                  string            `json:"path"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
}

type Response struct {
	Headers         map[string]string `json:"headers"`
	StatusCode      int               `json:"statusCode"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

func Handler(config *proxy.Configuration) func(*Request) (*Response, error) {
	return func(request *Request) (*Response, error) {
		requestURL, err := url.Parse(request.Path)
		if err != nil {
			return &Response{
				StatusCode: http.StatusBadRequest,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: err.Error(),
			}, nil
		}

		parameters := url.Values{}
		for k, v := range request.QueryStringParameters {
			parameters.Add(k, v)
		}
		requestURL.RawQuery = parameters.Encode()

		pr, err := proxy.NewRequestFromURL(requestURL)
		if err != nil {
			return &Response{
				StatusCode: http.StatusBadRequest,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: err.Error(),
			}, nil
		}

		rec := httptest.NewRecorder()
		proxy.Proxy(config, rec, pr)
		result := rec.Result()
		defer result.Body.Close()

		data, _ := ioutil.ReadAll(result.Body)

		headers := make(map[string]string)
		for k := range result.Header {
			headers[k] = result.Header.Get(k)
		}

		return &Response{
			StatusCode:      result.StatusCode,
			Body:            base64.StdEncoding.EncodeToString(data),
			Headers:         headers,
			IsBase64Encoded: true,
		}, nil
	}
}

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := flags.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			// Exit with no error if --help was given. This is used to test the build.
			os.Exit(0)
		}
		log.Fatal(err)
	}

	config := &proxy.Configuration{}
	config.LoadEnvironmentVariables()

	lambda.Start(Handler(config))
}
