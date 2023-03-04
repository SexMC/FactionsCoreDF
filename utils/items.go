package utils

import (
	"path"

	"github.com/linuxtf/dragonfly/libraries/config"
	"github.com/linuxtf/dragonfly/server/item"
)

var items = config.New(path.Join("configs", "items.json"))

// ItemStack parses and returns an Item Stack from the config
func ItemStack(key string) item.Stack {
	return items.ItemStack(key)
}

// Items parses and returns a slice of ItemStack from the config
func Items(key string) []item.Stack {
	return items.Items(key)
}
