package commands

import (
	"fmt"
	"sync"

	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/world"
)

type FClaimCmd struct {
	Claim cmd.SubCommand `cmd:"claim"`

	Radius cmd.Optional[int] `cmd:"radius"`
}

func (c FClaimCmd) Run(src cmd.Source, o *cmd.Output) {
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

	radius, ok := c.Radius.Load()

	if !ok {
		radius = 1
	}

	radiusLimit := int(utils.GetFactionConfig[float64]("radius_claim_limit"))

	if radius <= 0 || radius > radiusLimit {
		p.Message(utils.Message("invalid_claim_radius", fmt.Sprint(radiusLimit)))
		return
	}

	if radius == 1 {
		claimed := tryClaim(p, faction, chunkPos)

		if !claimed {
			owner := memory.ChunkOwner(chunkPos)

			p.Message(utils.Message("chunk_already_claimed", owner.Name))
			return
		}

		p.Message(utils.Message("chunk_claimed", chunkPos.X(), chunkPos.Z()))
		return
	}

	even := radius%2 == 0

	if even {
		radius--
	}

	start := int32((1 - radius) / 2)
	end := int32((radius - 1) / 2)
	claimedLands := 0
	var mu sync.Mutex

	// claim odd x by x claims
	var oddwg sync.WaitGroup

	for x := start; x < end+1; x++ {
		for z := start; z < end+1; z++ {
			oddwg.Add(1)

			go func(xc int32, zc int32) {
				chunk := &world.ChunkPos{chunkPos.X() + xc, chunkPos.Z() + zc}

				claimed := tryClaim(p, faction, chunk)
				if claimed {
					mu.Lock()
					claimedLands++
					mu.Unlock()
				}

				oddwg.Done()
			}(x, z)
		}
	}

	oddwg.Wait()

	if even {
		start = start - 1

		var evenwg1 sync.WaitGroup

		for z := start + 1; z < end+1; z++ {
			evenwg1.Add(1)

			go func(zc int32) {
				chunk := &world.ChunkPos{chunkPos.X() + start, chunkPos.Z() + zc}

				claimed := tryClaim(p, faction, chunk)
				if claimed {
					mu.Lock()
					claimedLands++
					mu.Unlock()
				}

				evenwg1.Done()
			}(z)
		}

		var evenwg2 sync.WaitGroup

		for x := end; x > start-1; x-- {
			evenwg2.Add(1)

			go func(xc int32) {
				chunk := &world.ChunkPos{chunkPos.X() + xc, chunkPos.Z() + start}

				claimed := tryClaim(p, faction, chunk)
				if claimed {
					mu.Lock()
					claimedLands++
					mu.Unlock()
				}

				evenwg2.Done()
			}(x)
		}

		evenwg2.Wait()
		evenwg1.Wait()

		radius++
	}

	if claimedLands == 0 {
		p.Message(utils.Message("radius_claim_failed", fmt.Sprint(radius, " x ", radius)))
		return
	}

	p.Message(utils.Message("radius_claim_successful", fmt.Sprint(claimedLands), fmt.Sprint(radius, " x ", radius)))
}

// tryClaim tries to claim the land, returns whether successful or not
func tryClaim(p *player.Player, faction *factions.Faction, chunkPos *world.ChunkPos) bool {
	if memory.ChunkOwner(chunkPos) != nil {
		return false
	}

	memory.RegisterClaim(faction, chunkPos)
	return true
}
