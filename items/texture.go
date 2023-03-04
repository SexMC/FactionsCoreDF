package items

import (
	"image"
	"os"

	"github.com/linuxtf/dragonfly/libraries/console"
)

// texture returns the texture of a custom item
func texture(id string) image.Image {
	imgFile, err := os.Open("./textures/" + id)

	if err != nil {
		console.Log.Error(err)
		return nil
	}

	img, _, err := image.Decode(imgFile)

	if err != nil {
		console.Log.Error(err)
		return nil
	}

	return img.(image.Image)
}
