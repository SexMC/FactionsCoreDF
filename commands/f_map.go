package commands

import (
	"fmt"

	"github.com/inceptionmc/factions/factions/board"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FMapCmd struct {
	Map cmd.SubCommand `cmd:"map"`

	AutoUpdate cmd.Optional[bool] `cmd:"autoupdate"`
}

func (c FMapCmd) Run(src cmd.Source, o *cmd.Output) {
	p, ok := src.(*player.Player)

	if !ok {
		o.Print(utils.Message("command_usage_by_console"))
		return
	}

	fPlayer := memory.FPlayer(p)

	autoUpdate, ok := c.AutoUpdate.Load()
	if !ok {
		autoUpdate = false
	}

	if fPlayer.AutoUpdate != autoUpdate {
		p.Message(utils.Message("faction_map_autoupdate", fmt.Sprint(autoUpdate)))
		fPlayer.AutoUpdate = autoUpdate
		return
	}

	// send the board
	p.Message(board.FactionMap(fPlayer))
}
