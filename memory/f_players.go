package memory

import (
	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/factions/chat"
	"github.com/inceptionmc/factions/postgres"
	"github.com/linuxtf/dragonfly/server/player"
)

var FPlayers = map[string]*factions.FPlayer{}

// LoadFPlayer loads the Faction player data from the database into the memory
func LoadFPlayer(p *player.Player) {
	if postgres.FPlayerExists(p.XUID()) {
		fPlayer, fac := postgres.FPlayer(p)
		fPlayer.Faction = Faction(fac)

		FPlayers[p.XUID()] = fPlayer
		return
	}

	NewFPlayer(p)
}

// NewFPlayer creates a new Faction player data
func NewFPlayer(p *player.Player) {
	fPlayer := &factions.FPlayer{
		Player:  p,
		Faction: nil,
		Channel: chat.GlobalChannel{},
	}

	FPlayers[p.XUID()] = fPlayer
	postgres.SaveFPlayer(fPlayer)
}

// FPlayer gets the Player Data for a player
func FPlayer(p *player.Player) *factions.FPlayer {
	data := FPlayers[p.XUID()]
	return data
}

// SaveFPlayer unloads the Faction player data from the memory and saves it to the database
func SaveFPlayer(xuid string) {
	data := FPlayers[xuid]
	postgres.UpdateFPlayer(data)

	delete(FPlayers, xuid)
}

// SaveAllFPlayers saves all the Faction Player data from the memory into the database
func SaveAllFPlayers() {
	for xuid := range FPlayers {
		SaveFPlayer(xuid)
	}
}
