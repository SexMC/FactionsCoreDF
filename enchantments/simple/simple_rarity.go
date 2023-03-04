package simple

import "github.com/inceptionmc/factions/utils"

type SimpleRarity struct{}

func (SimpleRarity) Name() string {
	return utils.EnchantmentRarityName("simple")
}

func (SimpleRarity) Cost() int {
	return utils.EnchantmentRarityCost("simple")
}

func (SimpleRarity) Weight() int {
	return utils.EnchantmentRarityWeight("simple")
}
