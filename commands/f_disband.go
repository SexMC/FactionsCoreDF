package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/ui"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FDisbandCmd struct {
	Disband cmd.SubCommand `cmd:"disband"`
}

func (FDisbandCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	fPlayer := memory.FPlayer(p)

	if !fPlayer.IsInAnyFaction() {
		p.Message(utils.Message("must_be_in_a_faction"))
		return
	} else if fPlayer.Faction.Leader.Xuid != p.XUID() {
		p.Message(utils.Message("must_be_leader"))
		return
	}

	p.SendForm(ui.NewFDisbandUI(fPlayer.Faction.Name))
}
