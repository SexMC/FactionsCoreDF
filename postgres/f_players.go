package postgres

import (
	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/factions/chat"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/libraries/srv/postgres"
	"github.com/linuxtf/dragonfly/server/player"
)

// FPlayer returns the data from the database
func FPlayer(p *player.Player) (fPlayer *factions.FPlayer, faction string) {
	rows := postgres.Query(`SELECT FACTION, CHANNEL FROM FPLAYERS WHERE XUID = $1`, p.XUID())

	if rows.Next() {
		var fac, channel string
		rows.Scan(&fac, &channel)

		fPlayer := &factions.FPlayer{
			Player:  p,
			Faction: nil,
			Channel: utils.ChannelFromID(channel),
		}
		return fPlayer, fac
	}

	return nil, ""
}

// SaveFPlayer saves the player data into the Database
func SaveFPlayer(fPlayer *factions.FPlayer) {
	faction := fPlayer.Faction
	var fName string

	if faction != nil {
		fName = faction.Name
	}

	postgres.Exec(`INSERT INTO FPLAYERS(XUID, FACTION, CHANNEL) VALUES($1, $2, $3)`, fPlayer.Player.XUID(), fName, chat.ChannelID(fPlayer.Channel))
}

// UpdateFPlayer updates the player data into the Database
func UpdateFPlayer(fPlayer *factions.FPlayer) {
	faction := fPlayer.Faction

	fName := ""
	if faction != nil {
		fName = faction.Name
	}

	postgres.Exec(`UPDATE FPLAYERS SET FACTION = $1, CHANNEL = $2 where XUID = $3`, fName, chat.ChannelID(fPlayer.Channel), fPlayer.Player.XUID())
}

// FPlayerExists returns whether the player data exists for a player
func FPlayerExists(xuid string) bool {
	rows := postgres.Query(`SELECT * FROM FPLAYERS WHERE XUID = $1`, xuid)
	return rows.Next()
}
