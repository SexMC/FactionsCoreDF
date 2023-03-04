package factions

import "github.com/inceptionmc/factions/vault"

// Faction Storage that contains all the faction's money, TNT, spawners, etc.
type Storage struct {
	// money is the map of deposited Faction Money
	money map[string]int32

	// tnt is the deposited Faction TNT
	tnt int32

	// spawners is the amount of spawners the Faction owns
	spawners map[int]int32

	// vaults are the faction vaults that are accessible by
	// the members of a faction.
	vaults map[int]*vault.Vault
}

// NewStorage creates and returns a newly initialised storage
func NewStorage() *Storage {
	return &Storage{
		money:    map[string]int32{},
		spawners: map[int]int32{},
		vaults:   map[int]*vault.Vault{},
	}
}

// Money is the total money owned by the faction
func (b Storage) Money() int32 {
	money := int32(0)

	for _, v := range b.money {
		money += v
	}

	return money
}

// TNT is the total TNT owned by the faction
func (b Storage) TNT() int32 {
	return b.tnt
}

// Worth is the total worth of the Faction
func (b Storage) Worth() int32 {
	worth := int32(0)

	for _, v := range b.money {
		worth += v
	}

	return worth
}
