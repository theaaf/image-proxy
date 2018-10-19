package proxy

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	assert.EqualValues(t, expectedWidth, img.Bounds().Dx())
	assert.EqualValues(t, expectedHeight, img.Bounds().Dy())
}

func testScalingFilter(t *testing.T, filter Filter, imageName string, expectedWidth, expectedHeight int) {
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
			testScalingFilter(t, filter, "50x50.png", 50, 50)
		})

		t.Run("Tall", func(t *testing.T) {
			testScalingFilter(t, filter, "50x400.png", 50, 400)
		})

		t.Run("TallAndWide", func(t *testing.T) {
			testScalingFilter(t, filter, "400x400.png", 200, 200)
		})

		t.Run("Wide", func(t *testing.T) {
			testScalingFilter(t, filter, "200x200.png", 200, 200)
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
					testScalingFilter(t, filter, "50x50.png", 50, 50)
				})

				t.Run("Tall", func(t *testing.T) {
					testScalingFilter(t, filter, "50x400.png", 50, 200)
				})

				t.Run("TallAndWide", func(t *testing.T) {
					testScalingFilter(t, filter, "400x400.png", 100, 200)
				})

				t.Run("Wide", func(t *testing.T) {
					testScalingFilter(t, filter, "200x200.png", 100, 200)
				})
			})
		}
	})
}

func TestJPEGFilter(t *testing.T) {
	in := getResponse("50x50.png")
	out, filterErr := JPEGFilter(98)(in)
	require.Nil(t, filterErr)
	assert.Equal(t, "image/jpeg", out.Header.Get("Content-Type"))
	_, err := jpeg.Decode(out.Body)
	require.NoError(t, err)
}
