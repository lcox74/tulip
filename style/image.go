package style

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lcox74/tulip/constants"
)

// BlackWhiteImage takes an image and converts it to black and white as the 
// SHARP display only supports black and white.
func BlackWhiteImage(img *rl.Image) {
    for x := int32(0); x < img.Width; x++ {
		for y := int32(0); y < img.Height; y++ {
			c := rl.GetImageColor(*img, x, y)

			sum := uint16(c.R) + uint16(c.G) + uint16(c.B)
			if sum < 255/3 {
				rl.ImageDrawPixel(img, x, y, constants.COLOR_FG)
			} else {
				rl.ImageDrawPixel(img, x, y, constants.COLOR_BG)
			}

		}
	}
}