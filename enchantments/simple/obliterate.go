package simple

import (
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/item"
	"github.com/linuxtf/dragonfly/server/world"
)

type Obliterate struct{}

func (Obliterate) Name() string {
	return utils.EnchantmentName("obliterate")
}

func (Obliterate) MaxLevel() int {
	return utils.MaxEnchantmentLevel("obliterate")
}

func (Obliterate) Cost(level int) (int, int) {
	return utils.MinEnchantmentCost("obliterate"), utils.MaxEnchantmentCost("obliterate")
}

func (Obliterate) Rarity() item.EnchantmentRarity {
	return SimpleRarity{}
}

func (Obliterate) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return true
}

func (Obliterate) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Sword)
	return ok
}
