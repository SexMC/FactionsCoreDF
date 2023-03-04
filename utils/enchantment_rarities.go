package utils

import (
	"path"

	"github.com/linuxtf/dragonfly/libraries/config"
)

var enchantmentRarities = config.New(path.Join("configs", "enchantment_rarities.json"))

// EnchantmentRarityName returns the enchantment rarity name
func EnchantmentRarityName(key string) string {
	data := enchantmentRarities.Config(key)
	return data.String("name")
}

// EnchantmentRarityCost returns the enchantment rarity cost
func EnchantmentRarityCost(key string) int {
	data := enchantmentRarities.Config(key)
	return data.Int("cost")
}

// EnchantmentRarityWeight returns the enchantment rarity weight
func EnchantmentRarityWeight(key string) int {
	data := enchantmentRarities.Config(key)
	return data.Int("weight")
}
