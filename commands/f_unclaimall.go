package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FUnclaimAllCmd struct {
	UnclaimAll cmd.SubCommand `cmd:"unclaimall"`

	Confirmation cmd.Optional[string] `cmd:"confirmation"`
}

func (c FUnclaimAllCmd) Run(src cmd.Source, o *cmd.Output) {
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
	if !utils.RankHasPermission(rank, "unclaimall") {
		mustBeRank := utils.RankWithNativePermission("unclaimall")
		p.Message(utils.Message("must_be_" + mustBeRank))

		return
	}

	// confirmation
	conf, ok := c.Confirmation.Load()

	if !ok {
		p.Message(utils.Message("confirmation_error", "f unclaimall", faction.Name))
		return
	}

	if conf != faction.Name {
		p.Message(utils.Message("confirmation_does_not_match"))
		return
	}

	// unclaim all the claims
	for _, c := range memory.GetClaims(faction.Name) {
		memory.DeleteClaim(c.Position)
	}
}
