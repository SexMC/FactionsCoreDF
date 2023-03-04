package factions

import (
	"math"

	"github.com/inceptionmc/factions/factions/chat"
	"github.com/inceptionmc/factions/redis"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/server/block/cube"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/world"
)

type FPlayer struct {
	Player *player.Player

	Faction       *Faction
	Channel       chat.Channel
	Chunk         *world.ChunkPos
	AutoClaim     bool
	AutoUpdate    bool
	ErrorCooldown int64
}

var relationColors = map[int]string{
	Neutral: "§a",
	Ally:    "§a",
	Truce:   "§d",
	Enemy:   "§c",
}

// IsInAnyFaction returns if the player is in any faction
func (f FPlayer) IsInAnyFaction() bool {
	return f.Faction != nil
}

// IsInFaction returns if the player is in a specific faction
func (f FPlayer) IsInFaction(faction string) bool {
	return f.Faction.Name == faction
}

// GetFMember returns the player's FMember instance in a faction
func (f FPlayer) GetFMember() *FMember {
	for _, m := range f.Faction.Members {
		if m.Xuid == f.Player.XUID() {
			return m
		}
	}

	return nil
}

// Invite invites the player to a Faction if not invited already
func (f FPlayer) Invite(faction *Faction, invitedBy *player.Player) {
	if redis.CheckInvite(f.Player.XUID(), faction.Name) {
		invitedBy.Message(utils.Message("already_invited", f.Player.Name()))
		return
	}

	redis.Invite(f.Player.XUID(), faction.Name, invitedBy.Name())

	f.Player.SendToast(utils.ToastTitle("faction_invitation"), utils.ToastContent("faction_invitation", faction.Name, faction.Name))
	f.Player.Message(utils.Message("invitation_received"))
	invitedBy.Message(utils.Message("invited_successfully", f.Player.Name()))
}

// SwitchChannel switches the chat channel of a Faction Player
func (f *FPlayer) SwitchChannel() {
	switch f.Channel.ChannelType() {
	case chat.Global:
		f.Channel = chat.TrucesChannel{}
	case chat.Truces:
		f.Channel = chat.AlliesChannel{}
	case chat.Allies:
		f.Channel = chat.FactionChannel{}
	case chat.Faction:
		if f.GetFMember().Rank >= Manager {
			f.Channel = chat.ModeratorChannel{}
			return
		}

		f.Channel = chat.GlobalChannel{}
	case chat.Moderator:
		f.Channel = chat.GlobalChannel{}
	}
}

// SetChannel sets the channel to a specific channel for a Faction Player
func (f *FPlayer) SetChannel(channel chat.Channel) {
	f.Channel = channel
}

// JoinFaction makes the player join the faction
func (f *FPlayer) JoinFaction(faction *Faction) {
	f.Faction = faction
}

// LeaveFaction removes the Faction for the player
func (f *FPlayer) LeaveFaction() {
	f.Faction = nil
	f.AutoClaim = false
}

// RelationColor returns the relation color of the player with a Faction
func (f *FPlayer) RelationColor(faction *Faction) string {
	if f.Faction == nil {
		return relationColors[Neutral]
	}

	rel := f.Faction.RelationWith(faction)
	return relationColors[rel]
}

// Direction returns the facing direction of the player
func (f *FPlayer) Direction() cube.Direction {
	degrees := math.Mod(f.Player.Rotation().Yaw()-180, 360)
	if degrees < 0 {
		degrees += 360
	}

	if degrees >= 0 && degrees < 22.5 {
		return cube.North
	} else if degrees >= 22.5 && degrees < 67.5 {
		return cube.NorthEast
	} else if degrees >= 67.5 && degrees < 112.5 {
		return cube.East
	} else if degrees >= 112.5 && degrees < 157.5 {
		return cube.SouthEast
	} else if degrees >= 157.5 && degrees < 202.5 {
		return cube.South
	} else if degrees >= 202.5 && degrees < 247.5 {
		return cube.SouthWest
	} else if degrees >= 247.5 && degrees < 292.5 {
		return cube.West
	} else if degrees >= 292.5 && degrees < 337.5 {
		return cube.NorthWest
	} else if degrees >= 337.5 && degrees < 360 {
		return cube.North
	}

	return cube.North
}
