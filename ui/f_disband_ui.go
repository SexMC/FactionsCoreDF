package ui

import (
	"fmt"

	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/factions/chat"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/libraries/broadcast"
	"github.com/linuxtf/dragonfly/libraries/srv"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/player/form"
)

type FDisbandUI struct {
	Input form.Input
}

// NewDisbandUI ...
func NewFDisbandUI(faction string) form.Form {
	data := utils.GetUI("f_disband_ui")
	input := data["input"].(map[string]any)

	return form.New(
		FDisbandUI{
			Input: form.Input{
				Text:        fmt.Sprintf(input["text"].(string), faction),
				Default:     input["default"].(string),
				Placeholder: input["placeholder"].(string),
			},
		},
		data["title"].(string),
	)
}

// Submit ...
func (f FDisbandUI) Submit(submitter form.Submitter) {
	p := submitter.(*player.Player)
	name := f.Input.Value()

	fPlayer := memory.FPlayer(p)
	if fPlayer.Faction.Name != name {
		p.Message(utils.Message("faction_name_does_not_match"))
		return
	}

	p.Message(utils.Message("faction_disbanded", name))
	broadcast.BroadcastMessage(utils.Message("broadcast_faction_disbanded", p.Name(), name))

	for _, m := range fPlayer.Faction.Members {
		player, ok := srv.Srv.PlayerByXUID(m.Xuid)

		if !ok {
			continue
		}

		if m.Rank != factions.Leader {
			player.Message(utils.Message("broadcast_faction_members_disbanded"))
		}

		fPlayer := memory.FPlayer(player)
		fPlayer.LeaveFaction()

		if fPlayer.Channel.ChannelType() != chat.Global {
			fPlayer.SetChannel(chat.GlobalChannel{})
		}
	}

	memory.DeleteFaction(name)
}
