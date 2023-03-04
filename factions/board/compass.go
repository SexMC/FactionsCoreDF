package board

import (
	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/block/cube"
)

var Directions = map[cube.Direction]string{
	cube.North:     "N",
	cube.NorthEast: "/",
	cube.East:      "E",
	cube.SouthEast: "\\",
	cube.South:     "S",
	cube.SouthWest: "/",
	cube.West:      "W",
	cube.NorthWest: "\\",
}

const Color_Active = utils.Red
const Color_Inactive = utils.Gold

// Compass generates and returns the compass of the player
func Compass(fPlayer *factions.FPlayer) []string {
	rows := [][]cube.Direction{
		{cube.NorthWest, cube.North, cube.NorthEast},
		{cube.West, cube.East},
		{cube.SouthWest, cube.South, cube.SouthEast},
	}

	dir := fPlayer.Direction()

	res := []string{}

	for i, r := range rows {
		row := ""

		for i2, n := range r {
			d := Directions[n]

			if i == 1 && i2 == 1 {
				row += Color_Inactive + "+"
			}

			if n == dir {
				row += Color_Active + d
			} else {
				row += Color_Inactive + d
			}
		}

		res = append(res, row)
	}

	return res
}
