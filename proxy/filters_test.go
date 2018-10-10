package proxy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"image/png"
	"io/ioutil"
	"net/http"
	"testing"
)

func getImage(name string) []byte {
	data, err := ioutil.ReadFile("testdata/" + name)
	if err != nil {
		panic(err)
	}
	return data
}

func getResponse(imageName string) *Response {
	reader := bytes.NewReader(getImage(imageName))
	return &Response{
		Header: http.Header{
			"Content-Type": []string{"image/png"},
		},
		Body: reader,
	}
}

func assertResponseImageSize(t *testing.T, resp *Response, expectedWidth, expectedHeight int) {
	img, err := png.Decode(resp.Body)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, expectedWidth, img.Bounds().Dx())
	assert.EqualValues(t, expectedHeight, img.Bounds().Dy())
}

func testFilter(t *testing.T, filter Filter, imageName string, expectedWidth, expectedHeight int) {
	resp, err := filter(getResponse(imageName))
	assert.Nil(t, err)
	assertResponseImageSize(t, resp, expectedWidth, expectedHeight)
}

func TestScalingFilter(t *testing.T) {
	t.Run("NoCropping", func(t *testing.T) {
		opts := &ScalingOptions{
			Fill: &Dimensions{100, 200},
		}
		filter := ScalingFilter(opts)

		t.Run("Small", func(t *testing.T) {
			testFilter(t, filter, "50x50.png", 50, 50)
		})

		t.Run("Tall", func(t *testing.T) {
			testFilter(t, filter, "50x400.png", 50, 400)
		})

		t.Run("TallAndWide", func(t *testing.T) {
			testFilter(t, filter, "400x400.png", 200, 200)
		})

		t.Run("Wide", func(t *testing.T) {
			testFilter(t, filter, "200x200.png", 200, 200)
		})
	})

	t.Run("Cropping", func(t *testing.T) {
		// all crop types should result in an image matching the exact dimensions, unless the provided image is smaller
		for cropType := range cropTypeToAnchor {
			t.Run(cropType.CamelCase(), func(t *testing.T) {
				opts := &ScalingOptions{
					Fill: &Dimensions{100, 200},
					Crop: &cropType,
				}
				filter := ScalingFilter(opts)

				t.Run("Small", func(t *testing.T) {
					testFilter(t, filter, "50x50.png", 50, 50)
				})

				t.Run("Tall", func(t *testing.T) {
					testFilter(t, filter, "50x400.png", 50, 200)
				})

				t.Run("TallAndWide", func(t *testing.T) {
					testFilter(t, filter, "400x400.png", 100, 200)
				})

				t.Run("Wide", func(t *testing.T) {
					testFilter(t, filter, "200x200.png", 100, 200)
				})
			})
		}
	})
}
