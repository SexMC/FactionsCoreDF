package ui

import (
	"fmt"

	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/player/form"
)

var temporary = map[string]*factions.Faction{}

type FInfoUI struct{}

// NewFInfoUI
func NewFInfoUI(player *player.Player, faction *factions.Faction) form.Form {
	temporary[player.XUID()] = faction

	data := utils.GetUI("f_info_menu")
	bttns := data["buttons"].(map[string]any)

	f := form.NewMenu(FInfoUI{}, data["title"])

	buttons := []form.Button{}

	for _, b := range bttns {
		buttons = append(buttons, form.NewButton(b.(map[string]any)["text"].(string), b.(map[string]any)["image"].(string)))
	}

	return f.WithButtons(
		buttons...,
	)
}

// Submit ...
func (f FInfoUI) Submit(submitter form.Submitter, pressed form.Button) {
	p := submitter.(*player.Player)

	switch pressed.Text {
	case "General Information":
		p.SendForm(newInfo(temporary[p.XUID()], 0))
		return
	case "Members":
		p.SendForm(newInfo(temporary[p.XUID()], 1))
		return
	case "Allies":
		p.SendForm(newInfo(temporary[p.XUID()], 2))
		return
	}

	delete(temporary, p.XUID())
}

type Info struct{}

// newInfo ...
func newInfo(faction *factions.Faction, panelType int) form.Form {
	data := utils.GetUI("f_info_panel")

	f := form.NewMenu(Info{}, fmt.Sprintf(data["title"].(string), faction.Name))

	switch panelType {
	case 0:
		f = form.Menu.WithBody(f, faction.GeneralInformation())
	}

	return f.WithButtons(
		form.NewButton("Done", ""), form.NewButton("§4Go back", ""),
	)
}

// Submit ...
func (f Info) Submit(submitter form.Submitter, pressed form.Button) {
	p := submitter.(*player.Player)

	switch pressed.Text {
	case "§4Go back":
		p.SendForm(NewFInfoUI(p, temporary[p.XUID()]))
	}
}
