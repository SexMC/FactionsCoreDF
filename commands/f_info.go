package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/ui"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FInfoCmd struct {
	Info cmd.SubCommand `cmd:"info"`

	Faction cmd.Optional[string] `cmd:"faction"`
}

func (c FInfoCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	fac, ok := c.Faction.Load()
	fPlayer := memory.FPlayer(p)
	faction := memory.Faction(fac)

	if !ok {
		if fPlayer.Faction == nil {
			p.Message(utils.Message("faction_not_found"))
			return
		}

		faction = fPlayer.Faction
	}

	if faction == nil {
		p.Message(utils.Message("faction_not_found"))
		return
	}

	p.SendForm(ui.NewFInfoUI(p, faction))
}
