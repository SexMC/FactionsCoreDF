package utils

import (
	"path"

	"github.com/linuxtf/dragonfly/libraries/config"
)

var enchantments = config.New(path.Join("configs", "enchantments.json"))

// EnchantmentName returns the enchantment name
func EnchantmentName(key string) string {
	data := enchantments.Config(key)
	return data.String("name")
}

// EnchantmentDescription returns the enchantment description
func EnchantmentDescription(key string) string {
	data := enchantments.Config(key)
	return data.String("description")
}

// MaxEnchantmentLevel returns the maximum enchantment level
func MaxEnchantmentLevel(key string) int {
	data := enchantments.Config(key)
	return data.Int("max_level")
}

// MinEnchantmentCost returns the minimum enchantment cost
func MinEnchantmentCost(key string) int {
	data := enchantments.Config(key)
	return data.Int("min_cost")
}

// EnchantmentCost returns the maximum enchantment cost
func MaxEnchantmentCost(key string) int {
	data := enchantments.Config(key)
	return data.Int("max_cost")
}

// EnchantmentProperty returns the custom enchantment property value
func EnchantmentProperty[T any](key string, property string) T {
	data := enchantments.Config(key)
	return data.Get(property).(T)
}
