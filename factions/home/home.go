package home

import (
	"github.com/inceptionmc/factions/factions/teleport"
	"github.com/linuxtf/dragonfly/server/block/cube"
	"github.com/linuxtf/dragonfly/server/player"
)

type Home struct {
	Location [3]float64

	Yaw float64

	Pitch float64

	Created int64

	SetBy string
}

func New(loc [3]float64, rotation cube.Rotation, created int64, setBy string) *Home {
	return &Home{
		Location: loc,
		Yaw:      rotation.Yaw(),
		Pitch:    rotation.Pitch(),
		Created:  created,
		SetBy:    setBy,
	}
}

func (h *Home) Teleport(p *player.Player) {
	teleport.Teleport(p, h.Location, h.Yaw, h.Pitch)
}
