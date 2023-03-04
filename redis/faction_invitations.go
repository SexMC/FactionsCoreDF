package redis

import (
	"strings"
	"time"

	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/libraries/console"
	"github.com/linuxtf/dragonfly/libraries/srv/redis"
)

// Invite sets the invite key for the player's XUID for the expiry duration
func Invite(xuid string, faction string, invitedBy string) {
	expiry := int(utils.GetFactionConfig[float64]("invitation_expiry"))

	err := redis.Conn.Set(redis.Ctx, xuid+":invitations:"+faction, invitedBy, time.Second*time.Duration(expiry)).Err()

	if err != nil {
		console.Log.Error(err)
	}
}

// AllInvites returns a map of all invites with the following data:
// Key: Faction name
// Value: Player from the faction that invited
func AllInvites(xuid string) map[string]string {
	invites := map[string]string{}
	keys, err := redis.Conn.Keys(redis.Ctx, xuid+":invitations:*").Result()

	if err != nil {
		console.Log.Error(err)
		return nil
	}

	for _, k := range keys {
		invitedBy, err := redis.Conn.Get(redis.Ctx, k).Result()

		if err != nil {
			console.Log.Error(err)
			return nil
		}

		faction := strings.Split(k, "invitations:")[1]
		invites[faction] = invitedBy
	}

	return invites
}

// InviteExpiry returns the expiry of a faction invite
func InviteExpiry(xuid string, faction string) time.Duration {
	expiry, err := redis.Conn.TTL(redis.Ctx, xuid+":invitations:"+faction).Result()

	if err != nil {
		console.Log.Error(err)
		return -1
	}

	return expiry
}

// CheckInvite returns whether the faction has already invited this player
func CheckInvite(xuid string, faction string) bool {
	invites := AllInvites(xuid)

	for f := range invites {
		if f == faction {
			return true
		}
	}

	return false
}

// DeleteInvite deletes the faction invite
func DeleteInvite(xuid string, faction string) {
	redis.Conn.Del(redis.Ctx, xuid+":invitations:"+faction)
}
