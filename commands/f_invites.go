package commands

import (
	"github.com/inceptionmc/factions/redis"
	"github.com/inceptionmc/factions/ui"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FInvitesCmd struct {
	Invite cmd.SubCommand `cmd:"invites"`
}

func (FInvitesCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	invites := redis.AllInvites(p.XUID())

	if len(invites) == 0 {
		p.Message(utils.Message("zero_invites"))
		return
	}

	p.SendForm(ui.NewFInvitesUI(p, invites))
}
