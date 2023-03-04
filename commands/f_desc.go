package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FDescCmd struct {
	Desc cmd.SubCommand `cmd:"desc"`

	Description string `cmd:"description"`
}

func (c FDescCmd) Run(src cmd.Source, o *cmd.Output) {
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

	faction := fPlayer.Faction
	fMember := fPlayer.GetFMember()
	rank := utils.RankID(fMember.Rank)

	// check if has permission
	if !utils.RankHasPermission(rank, "desc") {
		mustBeRank := utils.RankWithNativePermission("desc")
		p.Message(utils.Message("must_be_" + mustBeRank))

		return
	}

	// change the description
	faction.Description = c.Description
	p.Message(utils.Message("description_changed", faction.Description))
}
