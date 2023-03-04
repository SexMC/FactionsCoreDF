package commands

import (
	"time"

	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FAlertCmd struct {
	Alert cmd.SubCommand `cmd:"alert"`

	Type AlertType `cmd:"type"`

	Message string `cmd:"message"`
}

func (c FAlertCmd) Run(src cmd.Source, o *cmd.Output) {
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

	faction := fPlayer.Faction
	fMember := fPlayer.GetFMember()
	rank := utils.RankID(fMember.Rank)

	// check if has permission
	if !utils.RankHasPermission(rank, "alert") {
		mustBeRank := utils.RankWithNativePermission("alert")
		p.Message(utils.Message("must_be_" + mustBeRank))

		return
	}

	switch c.Type {
	case "title":
		titleData := utils.TitleData("faction_alert")
		fadeIn := time.Duration(titleData["fadeIn"].(float64)) * time.Second
		fadeOut := time.Duration(titleData["fadeOut"].(float64)) * time.Second
		stay := time.Duration(titleData["stay"].(float64)) * time.Second

		faction.BroadcastTitle(titleData["title"].(string), c.Message, fadeIn, fadeOut, stay)
	case "popup":
		faction.BroadcastPopup("§cFaction Alert: §f" + c.Message)
	case "notification":
		faction.BroadcastToast("§cFaction Alert §8[§7"+p.Name()+"§8]", c.Message)
	}
}

type AlertType string

func (AlertType) Type() string {
	return "type"
}

func (AlertType) Options(_ cmd.Source) []string {
	return []string{"title", "popup", "notification"}
}
