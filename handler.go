package main

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/inceptionmc/factions/factions"
	"github.com/inceptionmc/factions/factions/board"
	"github.com/inceptionmc/factions/factions/chat"
	"github.com/inceptionmc/factions/factions/teleport"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/utils"
	"github.com/linuxtf/dragonfly/libraries/broadcast"
	"github.com/linuxtf/dragonfly/server/block/cube"
	"github.com/linuxtf/dragonfly/server/event"
	"github.com/linuxtf/dragonfly/server/item"
	"github.com/linuxtf/dragonfly/server/player"
	"github.com/linuxtf/dragonfly/server/player/title"
	"github.com/linuxtf/dragonfly/server/world"
)

type PlayerHandler struct {
	player.NopHandler
	p *player.Player
}

// Handler ...
func Handler(p *player.Player) *PlayerHandler {
	return &PlayerHandler{p: p}
}

// HandleQuit ...
func (h *PlayerHandler) HandleQuit() {
	fPlayer := memory.FPlayer(h.p)
	faction := fPlayer.Faction

	if faction != nil && faction.OnlineCount() <= 1 {
		memory.SaveFaction(faction.Name)
	}

	memory.SaveFPlayer(h.p.XUID())

	// check if player was teleporting
	if teleport.IsTeleporting(h.p) {
		teleport.DeleteTeleportationData(h.p)
	}
}

// HandleMove ...
func (h *PlayerHandler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	// check if player was teleporting
	if teleport.IsTeleporting(h.p) {
		data := teleport.GetTeleportationData(h.p)

		// if position has changed, then delete the teleportation data
		if !data.CurrentPos.ApproxEqual(newPos) {
			teleport.DeleteTeleportationData(h.p)
		}
	}

	p := h.p
	fPlayer := memory.FPlayer(p)

	oldChunk := fPlayer.Chunk
	chunk := memory.ChunkPos(p.Position(), p.World())

	if fPlayer.Chunk == nil || oldChunk.X() != chunk.X() || oldChunk.Z() != chunk.Z() {
		fPlayer.Chunk = chunk
		owner := memory.ChunkOwner(chunk)

		// auto claim system
		if owner == nil && fPlayer.AutoClaim {
			memory.RegisterClaim(fPlayer.Faction, chunk)
			p.Message(fmt.Sprintf(utils.Message("chunk_claimed"), chunk.X(), chunk.Z()))
		}

		// faction map auto update system
		if fPlayer.AutoUpdate {
			p.Message(board.FactionMap(fPlayer))
		}

		// title system
		var oldOwner *factions.Faction

		if oldChunk == nil {
			return
		}
		oldOwner = memory.ChunkOwner(oldChunk)

		if oldOwner != nil && owner != nil && owner.Name == oldOwner.Name {
			return
		}

		if owner != nil {
			var relationColor string

			if fPlayer.Faction != nil {
				relationColor = fPlayer.RelationColor(owner)
			}

			p.SendTitle(title.New(relationColor + owner.Name).WithSubtitle("§7" + owner.Description).WithFadeInDuration(1 * time.Second).WithFadeOutDuration(1 * time.Second).WithDuration(1500 * time.Millisecond))
		}
	}
}

// HandleBlockBreak ...
func (h *PlayerHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	p := h.p
	w := p.World()

	chunk := memory.ChunkPos(pos.Vec3(), w)
	owner := memory.ChunkOwner(chunk)

	fPlayer := memory.FPlayer(p)

	if owner != nil && (fPlayer.Faction == nil || fPlayer.Faction.Name != owner.Name) {
		if time.Now().Unix() > fPlayer.ErrorCooldown {
			p.Message(utils.Message("cannot_modify_build", owner.Name))
			fPlayer.ErrorCooldown = time.Now().Unix() + 2
		}
		ctx.Cancel()
	}
}

// HandleBlockPlace ...
func (h *PlayerHandler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	p := h.p
	w := p.World()

	chunk := memory.ChunkPos(pos.Vec3(), w)
	owner := memory.ChunkOwner(chunk)

	fPlayer := memory.FPlayer(p)

	if owner != nil && (fPlayer.Faction == nil || fPlayer.Faction.Name != owner.Name) {
		if time.Now().Unix() > fPlayer.ErrorCooldown {
			p.Message(utils.Message("cannot_modify_build", owner.Name))
			fPlayer.ErrorCooldown = time.Now().Unix() + 2
		}
		ctx.Cancel()
	}
}

// HandleAttackEntity ...
func (h *PlayerHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	p := h.p
	t, ok := e.(*player.Player)

	if !ok {
		return
	}

	fPlayer := memory.FPlayer(p)
	tFPlayer := memory.FPlayer(t)

	owner := memory.ChunkOwner(memory.ChunkPos(t.Position(), t.World()))

	if tFPlayer.Faction != nil && fPlayer.Faction != nil {
		tFaction := tFPlayer.Faction
		faction := fPlayer.Faction

		if owner != nil && owner.Name != tFaction.Name && faction.Neutral(tFaction) {
			return
		}

		if faction.Enemy(tFaction) {
			return
		}

		if faction.Alliance(tFaction) {
			p.Message(utils.Message("cannot_hit_allies", t.Name()))
		}

		if faction.Truce(tFaction) {
			p.Message(utils.Message("cannot_hit_truces", t.Name()))
		}

		if faction.Neutral(tFaction) {
			p.Message(utils.Message("cannot_hit_neutrals", t.Name()))
		}
	}

	ctx.Cancel()
}

// HandleChat ...
func (h *PlayerHandler) HandleChat(ctx *event.Context, message *string) {
	player := h.p
	fPlayer := memory.FPlayer(player)

	channel := fPlayer.Channel
	format := utils.GetFactionConfig[map[string]any]("channels")[chat.ChannelID(channel)].(string)

	switch fPlayer.Channel.ChannelType() {
	case chat.Global:
		broadcast.BroadcastMessagef(fmt.Sprintf(format, player.Name(), *message))
	case chat.Truces:
		faction := fPlayer.Faction
		faction.Broadcast(fmt.Sprintf(format, player.Name(), *message))
	case chat.Allies:
		faction := fPlayer.Faction
		faction.Broadcast(fmt.Sprintf(format, player.Name(), *message))
	case chat.Faction:
		faction := fPlayer.Faction
		faction.Broadcast(fmt.Sprintf(format, player.Name(), *message))
	case chat.Moderator:
		faction := fPlayer.Faction
		faction.BroadcastManagers(fmt.Sprintf(format, player.Name(), *message))
	}
}
