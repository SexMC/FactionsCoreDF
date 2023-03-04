package commands

import (
	"fmt"

	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/ui"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FInviteCmd struct {
	Invite cmd.SubCommand `cmd:"invite"`

	Target cmd.Optional[[]cmd.Target] `cmd:"player"`
}

func (c FInviteCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	fPlayer := memory.FPlayer(p)

	if fPlayer.Faction == nil {
		p.Message(utils.Message("must_be_in_a_faction"))
		return
	} else if fPlayer.GetFMember().Rank < factions.Manager {
		p.Message(utils.Message("must_be_manager"))
		return
	}

	faction := fPlayer.Faction

	maxMembers := int(utils.GetFactionConfig[float64]("total_faction_members"))
	if maxMembers == len(faction.Members) {
		p.Message(utils.Message("max_members", fmt.Sprint(maxMembers)))
		return
	}

	targets, _ := c.Target.Load()

	if len(targets) == 0 {
		p.SendForm(ui.NewFInviteUI())
		return
	}

	if len(targets) > 1 {
		p.Message(utils.Message("more_than_one_target"))
		return
	}

	target := targets[0].(*player.Player)

	if target == p {
		p.Message(utils.Message("command_usage_on_self"))
		return
	}

	// check if target is in same faction
	if faction.IsMember(target) {
		p.Message(utils.Message("already_a_member", target.Name()))
		return
	}

	targetFPlayer := memory.FPlayer(target)
	targetFPlayer.Invite(faction, p)
}
