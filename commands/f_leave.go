package commands

import (
	"github.com/inceptionmc/factions/factions/chat"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FLeaveCmd struct {
	Leave cmd.SubCommand `cmd:"leave"`
}

func (FLeaveCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	fPlayer := memory.FPlayer(p)

	if !fPlayer.IsInAnyFaction() {
		p.Message(utils.Message("must_be_in_a_faction"))
		return
	} else if fPlayer.Faction.Leader.Xuid == p.XUID() {
		p.Message(utils.Message("leader_cannot_leave"))
		return
	}

	// leave the faction
	faction := fPlayer.Faction

	p.Message(utils.Message("left_faction", faction.Name))
	faction.Broadcast(utils.Message("broadcast_member_left_faction", p.Name(), utils.RankName(fPlayer.GetFMember().Rank)))

	fPlayer.LeaveFaction()

	if fPlayer.Channel.ChannelType() != chat.Global {
		fPlayer.SetChannel(chat.GlobalChannel{})
	}
}
