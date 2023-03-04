package utils

import (
	"path"

	"github.com/linuxtf/dragonfly/libraries/config"
)

var messages = config.New(path.Join("configs", "messages.json"))

// Message returns the formatted message
func Message(message string, args ...any) string {
	return messages.Message(message, args...)
}
