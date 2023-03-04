package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FKickCmd struct {
	Kick cmd.SubCommand `cmd:"kick"`

	Target []cmd.Target `cmd:"member"`
}

func (c FKickCmd) Run(src cmd.Source, o *cmd.Output) {
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
	faction := fPlayer.Faction

	if !utils.RankHasPermission(utils.RankID(fMember.Rank), "kick") {
		mustBeRank := utils.RankWithNativePermission("kick")
		p.Message(utils.Message("must_be_" + mustBeRank))
		return
	}

	if len(c.Target) > 1 {
		p.Message(utils.Message("more_than_one_target"))
		return
	}

	target := c.Target[0].(*player.Player)

	if target.XUID() == p.XUID() {
		p.Message(utils.Message("command_usage_on_self"))
		return
	}

	if !faction.IsMember(target) {
		p.Message(utils.Message("player_not_in_faction", target.Name()))
		return
	}

	targetFPlayer := memory.FPlayer(target)
	targetFMember := targetFPlayer.GetFMember()

	if !fMember.Compare(targetFMember) {
		p.Message(utils.Message("must_be_higher_in_hierarchy", target.Name()))
		return
	}

	// Kick the target
	faction.RemoveMember(target)
	targetFPlayer.LeaveFaction()

	p.Message(utils.Message("kicked_member", target.Name()))
	faction.Broadcast(utils.Message("broadcast_kicked_member", p.Name(), target.Name()))

	target.Message(utils.Message("kicked_from_faction", fPlayer.Faction.Name))
}
