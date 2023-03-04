package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FWarpCmd struct {
	Warp cmd.SubCommand `cmd:"warp"`

	Name string `cmd:"name"`
}

func (c FWarpCmd) Run(src cmd.Source, o *cmd.Output) {
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

	if !utils.RankHasPermission(utils.RankID(fMember.Rank), "warp") {
		mustBeRank := utils.RankWithNativePermission("warp")
		p.Message(utils.Message("must_be_" + mustBeRank))
		return
	}

	if !faction.WarpExists(c.Name) {
		p.Message(utils.Message("faction_warp_does_not_exist", c.Name))
		return
	}

	// teleport to the faction warp
	faction.Warps[c.Name].Teleport(p)
}
