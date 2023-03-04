package factions

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/inceptionmc/factions/factions/home"
	"github.com/inceptionmc/factions/factions/warp"
	"github.com/linuxtf/dragonfly/libraries/srv"
	"github.com/linuxtf/dragonfly/server/block/cube"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/player/title"
)

type Faction struct {
	// Name is the name of the faction
	Name string

	// Description is the description of the faction
	Description string

	// Allies are the slice of the names of the Allied factions
	Allies []string

	// Truces are the slices of the names of the truce factions
	Truces []string

	// Enemies are the slices of the names of the Enemy factions
	Enemies []string

	// Leader is the pointer to the faction leader
	Leader *FMember

	// Members is the slice of pointers to the faction members
	Members []*FMember

	// Home is the location where the faction members get teleported upon executed /f home
	Home *home.Home

	// Warps are the list of Faction Warps accessible by a Faction member
	Warps map[string]*warp.Warp

	// Storage is the storage for money, TNT, spawners and various other stuff that
	// the faction owns.
	Storage *Storage
}

const (
	Neutral = iota
	Ally
	Truce
	Enemy
)

// RelationWith returns the relation of a Faction with a Faction, does not return one sided relations
func (f *Faction) RelationWith(fac *Faction) int {
	if f.Alliance(fac) {
		return Ally
	}

	if f.Truce(fac) {
		return Truce
	}

	if f.Enemy(fac) {
		return Enemy
	}

	return Neutral
}

// Alliance returns whether two factions are in an alliance
func (f *Faction) Alliance(faction *Faction) bool {
	return f.MarkedAlly(faction) && faction.MarkedAlly(f)
}

// Truce returns whether two factions are in a truce relation
func (f *Faction) Truce(faction *Faction) bool {
	return f.MarkedTruce(faction) && faction.MarkedTruce(f)
}

// Enemy returns whether two factions are in a enemy relation
func (f *Faction) Enemy(faction *Faction) bool {
	return f.MarkedEnemy(faction) || faction.MarkedEnemy(f)
}

// Neutral returns whether two factions are in a neutral relation
func (f *Faction) Neutral(faction *Faction) bool {
	return f.MarkedNeutral(faction) && faction.MarkedNeutral(f)
}

// MarkAlly marks a faction as this faction's ally
func (f *Faction) MarkAlly(faction *Faction) {
	f.Allies = append(f.Allies, faction.Name)

	f.UnmarkEnemy(faction)
	f.UnmarkTruce(faction)
}

// UnmarkAlly unmarks a faction as this faction's ally
func (f *Faction) UnmarkAlly(faction *Faction) {
	for i, t := range f.Allies {
		if t == faction.Name {
			f.Allies = append(f.Allies[:i], f.Allies[i+1:]...)
			break
		}
	}
}

// MarkTruce marks a faction as this faction's truce
func (f *Faction) MarkTruce(faction *Faction) {
	f.Truces = append(f.Truces, faction.Name)

	f.UnmarkAlly(faction)
	f.UnmarkEnemy(faction)
}

// UnmarkTruce unmarks a faction as this faction's truce
func (f *Faction) UnmarkTruce(faction *Faction) {
	for i, t := range f.Truces {
		if t == faction.Name {
			f.Truces = append(f.Truces[:i], f.Truces[i+1:]...)
			break
		}
	}
}

// MarkEnemy marks a faction as this faction's enemy
func (f *Faction) MarkEnemy(faction *Faction) {
	f.Enemies = append(f.Enemies, faction.Name)

	f.UnmarkAlly(faction)
	f.UnmarkTruce(faction)
}

// UnmarkEnemy unmarks a faction as this faction's enemy
func (f *Faction) UnmarkEnemy(faction *Faction) {
	for i, t := range f.Enemies {
		if t == faction.Name {
			f.Enemies = append(f.Enemies[:i], f.Enemies[i+1:]...)
			break
		}
	}
}

// MarkNeutral marks a faction as this faction's neutral
func (f *Faction) MarkNeutral(faction *Faction) {
	f.UnmarkAlly(faction)
	f.UnmarkEnemy(faction)
	f.UnmarkTruce(faction)
}

// AddMember adds a new player to this faction
func (f *Faction) AddMember(player *player.Player) {
	f.Members = append(f.Members, &FMember{
		Name: player.Name(),
		Xuid: player.XUID(),
		Rank: Recruit,
	})
}

// RemoveMember removes a member from this faction
func (f *Faction) RemoveMember(player *player.Player) {
	for i, m := range f.Members {
		if m.Xuid == player.XUID() {
			f.Members = append(f.Members[:i], f.Members[i+1:]...)
			break
		}
	}
}

// IsMember returns whether a player is a part of this faction
func (f *Faction) IsMember(player *player.Player) bool {
	for _, m := range f.Members {
		return m.Xuid == player.XUID()
	}

	return false
}

// TryGetMember tries to returns the faction member if exists
func (f *Faction) TryGetMember(member string) *FMember {
	for _, m := range f.Members {
		if strings.EqualFold(m.Name, member) {
			return m
		}
	}

	return nil
}

// OnlineMembers returns the slice of online faction members
func (f *Faction) OnlineMembers() []*player.Player {
	members := []*player.Player{}

	for _, m := range f.Members {
		if p, ok := srv.Srv.PlayerByXUID(m.Xuid); ok {
			members = append(members, p)
		}
	}

	return members
}

// MembersWithRank returns the slice of faction members who are at a specific rank
func (f *Faction) MembersWithRank(rank int) []*FMember {
	members := []*FMember{}

	for _, m := range f.Members {
		if m.Rank == rank {
			members = append(members, m)
		}
	}

	return members
}

// OnlineCount returns the number of online faction members
func (f *Faction) OnlineCount() int {
	return len(f.OnlineMembers())
}

// GeneralInformation returns formatted text of the General Information of the Faction
func (f *Faction) GeneralInformation() string {
	info := []string{
		fmt.Sprintf("§7Faction Name: §a%s", f.Name),
		fmt.Sprintf("§7Leader: §a%s", f.Leader.Name),
		fmt.Sprintf("§7Online Members: §a%d", f.OnlineCount()),
	}

	return strings.Join(info, "\n")
}

// Broadcast broadcasts a message to all the Faction Members
func (f *Faction) Broadcast(message string) {
	for _, m := range f.OnlineMembers() {
		m.Message(message)
	}
}

// BroadcastManagers broadcasts a message to all the Faction Managers and higher
func (f *Faction) BroadcastManagers(message string) {
	for _, m := range f.Members {
		if p, ok := srv.Srv.PlayerByXUID(m.Xuid); ok {
			if m.Rank < Manager {
				continue
			}

			p.Message(message)
		}
	}
}

// BroadcastTitle sends a title on the screen of all the faction members
func (f *Faction) BroadcastTitle(t string, subtitle string, fadeIn, fadeOut, stay time.Duration) {
	title := title.New(t).WithSubtitle(subtitle).WithFadeInDuration(fadeIn).WithFadeOutDuration(fadeOut).WithDuration(stay)

	for _, m := range f.OnlineMembers() {
		m.SendTitle(title)
	}
}

// BroadcastPopup sends a action bar text on the screen of all the faction members
func (f *Faction) BroadcastPopup(msg string) {
	for _, m := range f.OnlineMembers() {
		m.SendPopup(msg)
	}
}

// BroadcastToast sends a toast on the screen of all faction members
func (f *Faction) BroadcastToast(title string, content string) {
	for _, m := range f.OnlineMembers() {
		m.SendToast(title, content)
	}
}

// SetWarp sets a warp at a location
func (f *Faction) SetWarp(name string, loc mgl64.Vec3, rotation cube.Rotation, createdBy string) {
	f.Warps[name] = warp.New(name, loc, rotation, time.Now().Unix(), createdBy)
}

// WarpExists returns whether a warp with a name exists
func (f *Faction) WarpExists(name string) bool {
	_, ok := f.Warps[name]
	return ok
}

// RemoveWarp removes a faction warp
func (f *Faction) RemoveWarp(name string) {
	delete(f.Warps, name)
}

// MarkedAlly returns one sided relation of a faction with a faction
func (f *Faction) MarkedAlly(faction *Faction) bool {
	for _, a := range f.Allies {
		if a == faction.Name {
			return true
		}
	}

	return false
}

// MarkedTruce returns one sided relation of a faction with a faction
func (f *Faction) MarkedTruce(faction *Faction) bool {
	for _, t := range f.Truces {
		if t == faction.Name {
			return true
		}
	}

	return false
}

// MarkedEnemy returns one sided relation of a faction with a faction
func (f *Faction) MarkedEnemy(faction *Faction) bool {
	for _, e := range f.Enemies {
		if e == faction.Name {
			return true
		}
	}

	return false
}

// MarkedNeutral returns one sided relation of a faction with a faction
func (f *Faction) MarkedNeutral(faction *Faction) bool {
	return !f.MarkedAlly(faction) && !f.MarkedEnemy(faction) && !f.MarkedTruce(faction)
}
