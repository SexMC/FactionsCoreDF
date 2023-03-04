package board

import (
	"fmt"
	"strings"

	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/world"
)

const MAP_KEY_CHARS = "\\/#?ç¬£$%=&^ABCDEFGHJKLMNOPQRSTUVWXYZÄÖÜÆØÅ1234567890abcdeghjmnopqrsuvwxyÿzäöüæøåâêîûô"

const MAP_WIDTH = int32(48)
const MAP_HEIGHT = int32(10)

const MAP_KEY_MIDDLE = utils.Aqua + "+"
const MAP_KEY_WILDERNESS = utils.Gray + "-"
const MAP_KEY_OVERFLOW = utils.White + "-" + utils.Reset

// FactionMap generates the Faction Map and returns it
func FactionMap(f *factions.FPlayer) string {
	chunk := f.Chunk
	centerFac := memory.ChunkOwner(chunk)

	compass := Compass(f)

	legend := map[rune]*factions.Faction{}
	characterIndex := 0
	overflown := false

	res := []string{}

	var area = "Wilderness"
	var color = utils.DarkGreen

	if centerFac != nil {
		area = centerFac.Name
		color = utils.Red
	}

	res = append(res, utils.Message("faction_map_header", color, fmt.Sprint(chunk.X(), ", ", chunk.Z()), area))

	for dz := int32(0); dz < MAP_HEIGHT; dz++ {
		row := ""

		for dx := int32(0); dx < MAP_WIDTH; dx++ {
			cx := chunk.X() - (MAP_WIDTH / 2) + dx
			cz := chunk.Z() - (MAP_HEIGHT / 2) + dz

			if chunk.X() == cx && chunk.Z() == cz {
				row += MAP_KEY_MIDDLE
				continue
			}

			fac := memory.ChunkOwner(&world.ChunkPos{cx, cz})

			var symbol rune

			if fac == nil {
				row += MAP_KEY_WILDERNESS
			} else {
				found := false

				for i, l := range legend {
					if l.Name == fac.Name {
						found = true
						symbol = i

						row += f.RelationColor(fac) + string(symbol)
						break
					}
				}

				if !found {
					if overflown {
						row += MAP_KEY_OVERFLOW
					} else {
						characterIndex++
						symbol = []rune(MAP_KEY_CHARS)[characterIndex]

						legend[symbol] = fac

						if characterIndex == len(MAP_KEY_CHARS) {
							overflown = true
						}

						row += f.RelationColor(fac) + string(symbol)
					}
				}
			}
		}

		if dz <= 2 {
			row = compass[dz] + utils.Substring(row, 9, len(row))
		}

		res = append(res, row)
	}

	var symbolMapping string

	for symbol, fac := range legend {
		color := f.RelationColor(fac)
		symbolMapping += color + string(symbol) + "§7: §f" + fac.Name + " "
	}

	res = append(res, symbolMapping)

	if overflown {
		res = append(res, utils.Message("faction_map_overflown"))
	}

	return strings.Join(res, "\n")
}
