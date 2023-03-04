package postgres

import (
	"encoding/json"
	"strings"

	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/factions/home"
	"github.com/inceptionmc/factions/factions/warp"
	"github.com/linuxtf/dragonfly/libraries/commons"
	"github.com/linuxtf/dragonfly/libraries/srv/postgres"
)

// Faction returns the Faction Data from the database
func Faction(faction string) *factions.Faction {
	rows := postgres.Query(`SELECT * FROM FACTIONS WHERE NAME = $1`, faction)

	if rows.Next() {
		var name, desc, a, t, e, l, m, h, w, s string
		rows.Scan(&name, &desc, &a, &t, &e, &l, &m, &h, &w)

		var allies, truces, enemies []string
		var leader *factions.FMember
		var members []*factions.FMember
		var home *home.Home
		var warps map[string]*warp.Warp
		var storage *factions.Storage

		json.Unmarshal([]byte(a), &allies)
		json.Unmarshal([]byte(t), &truces)
		json.Unmarshal([]byte(e), &enemies)
		json.Unmarshal([]byte(l), &leader)
		json.Unmarshal([]byte(m), &members)
		json.Unmarshal([]byte(h), &home)
		json.Unmarshal([]byte(w), &warps)
		json.Unmarshal([]byte(s), &storage)

		return &factions.Faction{
			Name:        name,
			Description: desc,
			Allies:      allies,
			Truces:      truces,
			Enemies:     enemies,
			Leader:      leader,
			Members:     members,
			Home:        home,
			Warps:       warps,
			Storage:     storage,
		}
	}

	return nil
}

// SaveFaction saves the faction data into the Database
func SaveFaction(faction *factions.Faction) {
	go func() {
		postgres.Exec(`INSERT INTO FACTIONS(NAME, DESCRIPTION, ALLIES, TRUCES, ENEMIES, LEADER, MEMBERS, HOME, WARPS, STORAGE) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`, faction.Name, faction.Description, commons.Encode(faction.Allies), commons.Encode(faction.Truces), commons.Encode(faction.Enemies), commons.Encode(faction.Leader), commons.Encode(faction.Members), commons.Encode(faction.Home), commons.Encode(faction.Warps), commons.Encode(faction.Storage))
	}()
}

// UpdateFaction saves the faction data into the Database
func UpdateFaction(faction *factions.Faction) {
	postgres.Exec(`UPDATE FACTIONS SET DESCRIPTION = $1, ALLIES = $2, TRUCES = $3, ENEMIES = $4, LEADER = $5, MEMBERS = $6, HOME = $7, WARPS = $8 WHERE NAME = $9`, faction.Description, commons.Encode(faction.Allies), commons.Encode(faction.Truces), commons.Encode(faction.Enemies), commons.Encode(faction.Leader), commons.Encode(faction.Members), commons.Encode(faction.Home), commons.Encode(faction.Warps), faction.Name)
}

// DeleteFaction clears the faction data from the database
func DeleteFaction(faction string) {
	postgres.Exec(`DELETE FROM FACTIONS WHERE NAME = $1`, faction)
}

// FactionExists returns whether a Faction with a name exists
func FactionExists(faction string) bool {
	rows := postgres.Query(`SELECT NAME FROM FACTIONS WHERE LOWER(NAME) = $1`, strings.ToLower(faction))
	return rows.Next()
}
