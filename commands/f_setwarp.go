package commands

import (
	"time"

	"github.com/inceptionmc/factions/factions/warp"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FSetWarpCmd struct {
	SetWarp cmd.SubCommand `cmd:"setwarp"`

	Name string `cmd:"name"`
}

func (c FSetWarpCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	fPlayer := memory.FPlayer(p)

	if fPlayer.Faction == nil {
		p.Message(utils.Message("must_be_in_a_faction"))
		return
	}

	fMember := fPlayer.GetFMember()
	faction := fPlayer.Faction

	if !utils.RankHasPermission(utils.RankID(fMember.Rank), "setwarp") {
		mustBeRank := utils.RankWithNativePermission("setwarp")
		p.Message(utils.Message("must_be_" + mustBeRank))
		return
	}

	// Check if warp exists
	if faction.WarpExists(c.Name) {
		p.Message(utils.Message("faction_warp_exists", c.Name))
		return
	}

	owner := memory.ChunkOwner(memory.ChunkPos(p.Position(), p.World()))

	// prevent setting home in chunk claimed by other faction
	if owner != nil && owner.Name != faction.Name {
		p.Message(utils.Message("cannot_set_warp", owner.Name))
		return
	}

	// Set the current location as the faction warp
	faction.Warps[c.Name] = warp.New(c.Name, p.Position(), p.Rotation(), time.Now().Unix(), p.Name())
	p.Message(utils.Message("faction_warp_set", c.Name))
}
