package koth

import (
	"github.com/linuxtf/dragonfly/server/block/cube"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/world"
)

// KothArena is the arena where the KoTH event takes place.
type KothArena struct {
	// Name is the unique ID / name of the arena
	Name string

	// Dimensions of the Koth Arena
	Dimensions cube.BBox

	// World in which the KothArena exists
	World *world.World
}

// Players returns the players that are inside the Koth arena
func (a *KothArena) Players(w *world.World) []*player.Player {
	players := []*player.Player{}

	entities := w.EntitiesWithin(a.Dimensions, func(e world.Entity) bool {
		_, ok := e.(*player.Player)
		return !ok
	})

	for _, e := range entities {
		players = append(players, e.(*player.Player))
	}

	return players
}
