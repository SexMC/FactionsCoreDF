package ui

import (
	"fmt"

	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/redis"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/player/form"
)

type FInvitesUI struct{}

var invitesCache = map[string]map[string]string{}
var factionCache = map[string]string{}

// NewFInvitesUI
func NewFInvitesUI(player *player.Player, invites map[string]string) form.Form {
	data := utils.GetUI("f_invites_menu")

	f := form.NewMenu(FInvitesUI{}, data["title"])
	f = f.WithBody(data["desc"])

	buttons := []form.Button{}

	for fac := range invites {
		buttons = append(buttons, form.NewButton(fmt.Sprintf(data["invite_title"].(string), fac), ""))
	}

	invitesCache[player.XUID()] = invites

	return f.WithButtons(
		buttons...,
	)
}

// Submit ...
func (f FInvitesUI) Submit(submitter form.Submitter, pressed form.Button) {
	p := submitter.(*player.Player)

	if pressed.Text == "" {
		delete(invitesCache, p.XUID())
		return
	}

	p.SendForm(newInvites(p, pressed.Text))
}

type Invites struct{}

// newInvites ...
func newInvites(player *player.Player, faction string) form.Form {
	data := utils.GetUI("f_invites_menu")

	f := form.NewMenu(Invites{}, data["title"])
	invites := invitesCache[player.XUID()]

	f = f.WithBody(fmt.Sprintf(data["invite_desc"].(string), faction, redis.InviteExpiry(player.XUID(), faction), invites[faction]))

	factionCache[player.XUID()] = faction

	return f.WithButtons(
		form.NewButton("Accept Invite", ""), form.NewButton("§4Reject Invite", ""),
	)
}

// Submit ...
func (f Invites) Submit(submitter form.Submitter, pressed form.Button) {
	p := submitter.(*player.Player)
	fPlayer := memory.FPlayer(p)

	switch pressed.Text {
	case "Accept Invite":
		if fPlayer.Faction != nil {
			p.Message(utils.Message("must_leave_faction", fPlayer.Faction.Name))
			return
		}

		// join the faction
		faction := memory.Faction(factionCache[p.XUID()])

		faction.Broadcast(utils.Message("broadcast_member_joined_faction", p.Name()))
		p.Message(utils.Message("joined_faction", factionCache[p.XUID()]))

		faction.AddMember(p)
		fPlayer.JoinFaction(faction)

		// delete the invite
		redis.DeleteInvite(p.XUID(), faction.Name)
		return
	case "§4Reject Invite":
		fac := factionCache[p.XUID()]
		p.Message(utils.Message("rejected_invite", fac))

		// delete the invite
		redis.DeleteInvite(p.XUID(), fac)

		// broadcast
		faction := memory.Faction(fac)
		faction.BroadcastManagers(utils.Message("broadcast_player_rejected_invite", p.Name()))

		delete(factionCache, fac)
		return
	}

	p.SendForm(NewFInvitesUI(p, invitesCache[p.XUID()]))
}
