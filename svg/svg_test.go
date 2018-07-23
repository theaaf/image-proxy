package svg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validSVG = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="100%" height="100%" viewBox="0 0 800 800" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:space="preserve" xmlns:serif="http://www.serif.com/" style="fill-rule:evenodd;clip-rule:evenodd;stroke-linejoin:round;stroke-miterlimit:1.41421;">
    <rect x="0" y="0" width="800" height="800" style="fill:url(#_Linear1);"/>
    <g transform="matrix(1,0,0,1,37.6063,17.7884)">
        <path d="M362.394,167.895L415.591,331.618L587.739,331.618L448.468,432.805L501.665,596.528L362.394,495.341L223.122,596.528L276.319,432.805L137.048,331.618L309.197,331.618L362.394,167.895Z" style="fill:white;"/>
    </g>
    <defs>
        <linearGradient id="_Linear1" x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(800,0,0,800,0,400)"><stop offset="0" style="stop-color:{{index . 0}};stop-opacity:1"/><stop offset="1" style="stop-color:{{index . 1}};stop-opacity:1"/></linearGradient>
    </defs>
</svg>`

func TestNew(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		svg, err := New([]byte(validSVG))
		assert.NoError(t, err)
		assert.NotNil(t, svg)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		svg, err := New([]byte("foo"))
		require.Error(t, err)
		assert.NotEmpty(t, err.Error())
		assert.Nil(t, svg)
	})
}
