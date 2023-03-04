package commands

import (
	"time"

	"github.com/inceptionmc/factions/factions/home"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FSetHomeCmd struct {
	SetHome cmd.SubCommand `cmd:"sethome"`
}

func (FSetHomeCmd) Run(src cmd.Source, o *cmd.Output) {
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

	if !utils.RankHasPermission(utils.RankID(fMember.Rank), "sethome") {
		mustBeRank := utils.RankWithNativePermission("sethome")
		p.Message(utils.Message("must_be_" + mustBeRank))
		return
	}

	// Set the current location as the faction home
	owner := memory.ChunkOwner(memory.ChunkPos(p.Position(), p.World()))

	// prevent setting home in chunk claimed by other faction
	if owner != nil && owner.Name != faction.Name {
		p.Message(utils.Message("cannot_set_home", owner.Name))
		return
	}

	faction.Home = home.New(p.Position(), p.Rotation(), time.Now().Unix(), p.Name())
	p.Message(utils.Message("faction_home_set"))
}
