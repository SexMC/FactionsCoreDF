package utils

import (
	"fmt"
	"path"

	"github.com/linuxtf/dragonfly/libraries/config"
)

var toasts = config.New(path.Join("configs", "toasts.json"))

// GetToast parses and returns a toast from a config file.
func GetToast(key string) (title string, content string) {
	return toasts.Toast(key)
}

// ToastTitle parses and returns the toast title from the config
func ToastTitle(key string, args ...any) string {
	title, _ := GetToast(key)
	return fmt.Sprintf(title, args...)
}

// ToastTitle parses and returns the toast title from the config
func ToastContent(key string, args ...any) string {
	_, content := GetToast(key)
	return fmt.Sprintf(content, args...)
}
