package commands

import (
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

type FUnclaimCmd struct {
	Unclaim cmd.SubCommand `cmd:"unclaim"`
}

func (FUnclaimCmd) Run(src cmd.Source, o *cmd.Output) {
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
	if !utils.RankHasPermission(rank, "claim") {
		mustBeRank := utils.RankWithNativePermission("claim")
		p.Message(utils.Message("must_be_" + mustBeRank))

		return
	}

	chunkPos := memory.ChunkPos(p.Position(), p.World())
	owner := memory.ChunkOwner(chunkPos)

	// check if land is already free
	if owner == nil {
		p.Message(utils.Message("chunk_is_wilderness"))
		return
	}

	// check if land is owned by another faction
	if owner != nil && owner.Name != faction.Name {
		p.Message(utils.Message("cannot_unclaim_chunk", owner.Name))
		return
	}

	// check if any warps or homes exist in this chunk
	// if they do, remove it
	if faction.Home != nil {
		chunk := memory.ChunkPos(faction.Home.Location, p.World())
		if chunk.X() == chunkPos.X() && chunk.Z() == chunkPos.Z() {
			faction.Home = nil
		}
	}

	for _, w := range faction.Warps {
		chunk := memory.ChunkPos(w.Location, p.World())
		if chunk.X() == chunkPos.X() && chunk.Z() == chunkPos.Z() {
			faction.RemoveWarp(w.Name)
		}
	}

	// unclaim
	memory.DeleteClaim(chunkPos)
	p.Message(utils.Message("chunk_unclaimed", chunkPos.X(), chunkPos.Z()))
}
