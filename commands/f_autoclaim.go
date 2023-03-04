package commands

import (
	"fmt"

	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FAutoClaimCmd struct {
	AutoClaim cmd.SubCommand `cmd:"autoclaim"`
}

func (FAutoClaimCmd) Run(src cmd.Source, o *cmd.Output) {
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
	rank := utils.RankID(fMember.Rank)

	// check if has permission
	if !utils.RankHasPermission(rank, "autoclaim") {
		mustBeRank := utils.RankWithNativePermission("autoclaim")
		p.Message(utils.Message("must_be_" + mustBeRank))

		return
	}

	fPlayer.AutoClaim = !fPlayer.AutoClaim
	p.Message(fmt.Sprintf(utils.Message("autoclaim_changed"), fPlayer.AutoClaim))
}
