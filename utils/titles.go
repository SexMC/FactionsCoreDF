package utils

import (
	"path"

	"github.com/linuxtf/dragonfly/libraries/config"
)

var titles = config.New(path.Join("configs", "titles.json"))

// TitleData returns the title data
func TitleData(key string) map[string]any {
	return titles.Map(key)
}
