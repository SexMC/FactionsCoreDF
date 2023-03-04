package memory

import (
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/postgres"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/libraries/commons"
	"github.com/linuxtf/dragonfly/server/world"
)

var Claims = sync.Map{}

// LoadClaims loads all the claims into the memory
func LoadClaims() int {
	for o, c := range postgres.GetAllClaims() {
		Claims.Store(o, c)
	}
	return utils.LenSyncMap(&Claims)
}

// GetClaims gets all the claims of a faction
func GetClaims(faction string) []*factions.Claim {
	claims := []*factions.Claim{}

	Claims.Range(func(key, value any) bool {
		c := value.(*factions.Claim)
		if c.Owner == faction {
			claims = append(claims, c)
		}

		return true
	})

	return claims
}

// RegisterClaim registers a new faction claim
func RegisterClaim(faction *factions.Faction, chunk *world.ChunkPos) {
	c := &factions.Claim{
		Position: chunk,
		Owner:    faction.Name,
		Created:  time.Now().Unix(),
	}

	Claims.Store(*chunk, c)
	postgres.RegisterClaim(c.Position, c.Owner, c.Created)
}

// DeleteClaim deletes a claim at a chunk
func DeleteClaim(chunk *world.ChunkPos) {
	Claims.Delete(*chunk)
	postgres.DeleteClaim(commons.Encode(chunk))
}

// FactionAt returns the faction that owns the chunk which contains a position
func FactionAt(pos mgl64.Vec3, w *world.World) *factions.Faction {
	chunk := ChunkPos(pos, w)
	return ChunkOwner(chunk)
}

// ChunkOwner returns the Faction that owns the chunk
func ChunkOwner(chunk *world.ChunkPos) *factions.Faction {
	data, ok := Claims.Load(*chunk)

	if ok {
		claim := data.(*factions.Claim)
		return Faction(claim.Owner)
	}

	return nil
}

// ChunkPos returns the chunk position that contains a position
func ChunkPos(pos mgl64.Vec3, w *world.World) *world.ChunkPos {
	chunkX := int(pos.X()) >> 4
	chunkZ := int(pos.Z()) >> 4
	return &world.ChunkPos{int32(chunkX), int32(chunkZ)}
}
