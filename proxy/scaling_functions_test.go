package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScaleToFit(t *testing.T) {
	t.Run("Small", func(t *testing.T) {
		w, h := ScaleToFit(100, 200)(50, 50)
		assert.Equal(t, w, 50)
		assert.Equal(t, h, 50)
	})

	t.Run("Tall", func(t *testing.T) {
		w, h := ScaleToFit(100, 200)(50, 400)
		assert.Equal(t, w, 25)
		assert.Equal(t, h, 200)
	})

	t.Run("Wide", func(t *testing.T) {
		w, h := ScaleToFit(100, 200)(200, 200)
		assert.Equal(t, w, 100)
		assert.Equal(t, h, 100)
	})
}

func TestScaleToFill(t *testing.T) {
	t.Run("Small", func(t *testing.T) {
		w, h := ScaleToFill(100, 200)(50, 50)
		assert.Equal(t, w, 50)
		assert.Equal(t, h, 50)
	})

	t.Run("Tall", func(t *testing.T) {
		w, h := ScaleToFill(100, 200)(50, 400)
		assert.Equal(t, w, 50)
		assert.Equal(t, h, 400)
	})

	t.Run("TallAndWide", func(t *testing.T) {
		w, h := ScaleToFill(100, 200)(400, 400)
		assert.Equal(t, w, 200)
		assert.Equal(t, h, 200)
	})

	t.Run("Wide", func(t *testing.T) {
		w, h := ScaleToFill(100, 200)(200, 200)
		assert.Equal(t, w, 200)
		assert.Equal(t, h, 200)
	})
}
