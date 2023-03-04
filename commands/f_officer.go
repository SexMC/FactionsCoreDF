package commands

import (
	"fmt"
	"strings"

	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/libraries/srv/users"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FOfficerCmd struct {
	Officer cmd.SubCommand `cmd:"officer"`

	Member string `cmd:"member"`
}

func (c FOfficerCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	// check if console
	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	fPlayer := memory.FPlayer(p)

	// check if faction exists
	if fPlayer.Faction == nil {
		p.Message(utils.Message("must_be_in_a_faction"))
		return
	}

	faction := fPlayer.Faction
	fMember := fPlayer.GetFMember()
	rank := utils.RankID(fMember.Rank)
	newRank := factions.Officer

	// check if has permission
	if !utils.RankHasPermission(rank, "officer") {
		mustBeRank := utils.RankWithNativePermission("officer")
		p.Message(utils.Message("must_be_" + mustBeRank))

		return
	}

	// check limits
	officersMax := int(utils.GetFactionConfig[float64]("number_of_officers"))
	if len(faction.MembersWithRank(newRank))+1 == officersMax {
		p.Message(utils.Message("max_officers", fmt.Sprint(officersMax)))
		return
	}

	// check target
	user := users.GetUserByName(c.Member)

	if !ok {
		p.Message(utils.Message("invalid_player"))
		return
	}

	if strings.EqualFold(user.Username, p.Name()) {
		p.Message(utils.Message("command_usage_on_self"))
		return
	}

	targetFMember := faction.TryGetMember(c.Member)

	if targetFMember == nil {
		p.Message(utils.Message("player_not_in_faction", user.Username))
		return
	}

	// set the rank
	targetFMember.Rank = newRank
	rank = utils.RankName(newRank)

	p.Message(utils.Message("faction_rank_changed", user.Username, rank))
	faction.Broadcast(utils.Message("broadcast_faction_rank_changed", p.Name(), user.Username, rank))
}
