package commands

import (
	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/factions/chat"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FChatCmd struct {
	Chat cmd.SubCommand `cmd:"chat"`

	Channel cmd.Optional[ChannelType] `cmd:"channel"`
}

func (c FChatCmd) Run(src cmd.Source, o *cmd.Output) {
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

	ch, ok := c.Channel.Load()

	if !ok {
		fPlayer.SwitchChannel()
		p.Message(utils.Message("chat_channel_changed", chat.ChannelName(fPlayer.Channel)))
		return
	}

	channel := utils.ChannelFromID(string(ch))

	if (channel == chat.ModeratorChannel{}) && fPlayer.GetFMember().Rank < factions.Manager {
		p.Message(utils.Message("must_be_manager"))
		return
	}

	fPlayer.SetChannel(channel)
	p.Message(utils.Message("chat_channel_changed", chat.ChannelName(fPlayer.Channel)))
}

type ChannelType string

func (ChannelType) Type() string {
	return "channel"
}

func (ChannelType) Options(_ cmd.Source) []string {
	return []string{"g", "global", "t", "truces", "a", "allies", "f", "faction", "m", "moderator"}
}
