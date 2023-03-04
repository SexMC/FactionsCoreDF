package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/ui"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FCreateCmd struct {
	Create cmd.SubCommand `cmd:"create"`
}

func (FCreateCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	fPlayer := memory.FPlayer(p)

	if fPlayer.Faction != nil {
		faction := fPlayer.Faction

		if faction.Leader.Xuid == p.XUID() {
			p.Message(utils.Message("must_disband_faction", faction.Name))
			return
		}

		p.Message(utils.Message("must_leave_faction", faction.Name))
		return
	}

	p.SendForm(ui.NewFCreateUI())
}
