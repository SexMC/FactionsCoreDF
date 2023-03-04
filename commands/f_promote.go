package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/libraries/srv"
	"github.com/linuxtf/dragonfly/libraries/srv/users"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/player/title"
)

type FPromoteCmd struct {
	Promote cmd.SubCommand `cmd:"promote"`

	Member string `cmd:"member"`
}

func (c FPromoteCmd) Run(src cmd.Source, o *cmd.Output) {
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

	// check if has permission
	if !utils.RankHasPermission(rank, "promote") {
		mustBeRank := utils.RankWithNativePermission("promote")
		p.Message(utils.Message("must_be_" + mustBeRank))

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

	// check if higher in hierarchy
	if !fMember.Compare(targetFMember) {
		p.Message(utils.Message("must_be_higher_in_hierarchy", user.Username))
		return
	}

	// check if rank is validated
	newRank := targetFMember.Successor()
	officersMax := int(utils.GetFactionConfig[float64]("number_of_officers"))
	managersMax := int(utils.GetFactionConfig[float64]("number_of_managers"))
	coleadersMax := int(utils.GetFactionConfig[float64]("number_of_coleaders"))

	switch newRank {
	case factions.Officer:
		if len(faction.MembersWithRank(newRank))+1 == officersMax {
			p.Message(utils.Message("max_officers", fmt.Sprint(officersMax)))
			return
		}
	case factions.Manager:
		if len(faction.MembersWithRank(newRank))+1 == managersMax {
			p.Message(utils.Message("max_managers", fmt.Sprint(managersMax)))
			return
		}
	case factions.CoLeader:
		if len(faction.MembersWithRank(newRank))+1 == coleadersMax {
			p.Message(utils.Message("max_coleaders", fmt.Sprint(coleadersMax)))
			return
		}
	case factions.Leader:
		p.Message(utils.Message("cannot_promote_to_leader"))
		return
	}

	// promote
	fMember.Rank = newRank
	faction.Broadcast(utils.Message("faction_member_promoted", targetFMember.Name, utils.RankName(newRank)))

	// send title
	targetPlayer, ok := srv.Srv.PlayerByName(user.Username)

	if !ok {
		return
	}

	titleData := utils.TitleData("promoted")
	fadeIn := time.Duration(titleData["fadeIn"].(float64)) * time.Second
	fadeOut := time.Duration(titleData["fadeOut"].(float64)) * time.Second
	stay := time.Duration(titleData["stay"].(float64)) * time.Second

	title := title.New(titleData["title"]).WithSubtitle(fmt.Sprintf(titleData["subtitle"].(string), utils.RankName(newRank))).WithFadeInDuration(fadeIn).WithFadeOutDuration(fadeOut).WithDuration(stay)
	targetPlayer.SendTitle(title)
}
