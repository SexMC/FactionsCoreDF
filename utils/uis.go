package utils

import (
	"path"

	"github.com/linuxtf/dragonfly/libraries/config"
)

var uis = config.New(path.Join("configs", "uis.json"))

// GetUI returns the UI data
func GetUI(ui string) map[string]any {
	return uis.Map(ui)
}
