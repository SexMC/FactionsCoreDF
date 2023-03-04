package simple

import (
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/item"
	"github.com/linuxtf/dragonfly/server/world"
)

type Headless struct{}

func (Headless) Name() string {
	return utils.EnchantmentName("headless")
}

func (Headless) MaxLevel() int {
	return utils.MaxEnchantmentLevel("headless")
}

func (Headless) Cost(level int) (int, int) {
	return utils.MinEnchantmentCost("headless"), utils.MaxEnchantmentCost("headless")
}

func (Headless) Rarity() item.EnchantmentRarity {
	return SimpleRarity{}
}

func (Headless) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return true
}

func (Headless) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Sword)
	return ok
}
