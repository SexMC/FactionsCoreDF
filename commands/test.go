package commands

import (
	"math"

	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// TestCmd is a test command, used for testing various new things
// temporary!!
type TestCmd struct{}

func (TestCmd) Run(src cmd.Source, o *cmd.Output) {
	p := src.(*player.Player)

	x := int32(math.Round(p.Position().X()))
	y := int32(math.Round(p.Position().Y()))
	z := int32(math.Round(p.Position().Z()))
	pos := protocol.BlockPos{x, y, z}

	//p.World().SetBlock(cube.PosFromVec3(p.Position()), block.Hopper{}, nil)

	wID := p.Session().NextWindowID()

	p.SendPacket(&packet.ContainerOpen{
		WindowID:                wID,
		ContainerPosition:       pos,
		ContainerType:           protocol.ContainerTypeHopper,
		ContainerEntityUniqueID: -1,
	})
}
