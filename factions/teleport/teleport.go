package teleport

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/player/title"
	"github.com/linuxtf/dragonfly/server/world/particle"
)

// TeleportationData contains the player undergoing teleportation's starting position, rotation
// and their target position
type TeleportationData struct {
	CurrentPos mgl64.Vec3

	TargetPos mgl64.Vec3

	TargetYaw float64

	TargetPitch float64
}

var teleportations = map[string]*TeleportationData{}

// Teleport teleports the player to a specific location
func Teleport(player *player.Player, loc mgl64.Vec3, yaw float64, pitch float64) {
	teleportations[player.XUID()] = &TeleportationData{
		CurrentPos:  player.Position(),
		TargetPos:   loc,
		TargetYaw:   yaw,
		TargetPitch: pitch,
	}

	player.Message(utils.Message("teleportation_started"))
	duration := int(utils.GetFactionConfig[float64]("teleport_duration"))

	go func() {
		for duration > -1 {
			if duration == 0 {
				data := utils.TitleData("teleported")
				t := data["title"].(string)

				// delete teleportation data
				DeleteTeleportationData(player)

				// teleport the player
				player.Teleport(loc)
				sendTitle(player, t, data)

				// particle effect
				player.World().AddParticle(player.Position(), particle.EndermanTeleport{})
				return
			} else {
				// check if teleportation was cancelled due to change in movement
				if !IsTeleporting(player) {
					data := utils.TitleData("teleportation_cancelled")
					t := data["title"].(string)

					// teleportation failure
					player.Message(utils.Message("teleportation_cancelled"))
					sendTitle(player, t, data)

					return
				} else {
					data := utils.TitleData("teleporting")
					t := fmt.Sprintf(data["title"].(string), duration)

					sendTitle(player, t, data)
				}
			}

			duration--
			time.Sleep(time.Second)
		}
	}()
}

// IsTeleporting returns whether a player is currently under the process of teleporting
func IsTeleporting(player *player.Player) bool {
	_, ok := teleportations[player.XUID()]
	return ok
}

// GetTeleportationData returns the TeleportationData of the player undergoing teleportation
func GetTeleportationData(player *player.Player) *TeleportationData {
	return teleportations[player.XUID()]
}

// DeleteTeleportationData deletes the teleportation data of the player
func DeleteTeleportationData(player *player.Player) {
	delete(teleportations, player.XUID())
}

// sendTitle sends a title to the player who is teleporting
func sendTitle(player *player.Player, t string, data map[string]any) {
	// build the title
	fadeIn := time.Duration(data["fadeIn"].(float64)) * time.Second
	fadeOut := time.Duration(data["fadeOut"].(float64)) * time.Second
	stay := time.Duration(data["stay"].(float64)) * time.Second

	// send the title
	title := title.New(t).WithSubtitle(data["subtitle"].(string)).WithFadeInDuration(fadeIn).WithFadeOutDuration(fadeOut).WithDuration(stay)
	player.SendTitle(title)
}
