package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FHomeCmd struct {
	Home cmd.SubCommand `cmd:"home"`
}

func (FHomeCmd) Run(src cmd.Source, o *cmd.Output) {
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

	if !utils.RankHasPermission(utils.RankID(fMember.Rank), "home") {
		mustBeRank := utils.RankWithNativePermission("home")
		p.Message(utils.Message("must_be_" + mustBeRank))
		return
	}

	if faction.Home == nil {
		p.Message(utils.Message("faction_has_no_home"))
		return
	}

	// teleport to the faction home
	faction.Home.Teleport(p)
}
