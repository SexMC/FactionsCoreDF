package vault

import (
	"github.com/linuxtf/dragonfly/server/item"
	"github.com/linuxtf/dragonfly/server/player"
)

// Vault is an e-chest that works the same way as enderchest does.
// However, unlike enderchests, e-chests can be made to access by one or more
// entities.
type Vault struct {
	storage map[int]item.Stack
}

// Open opens a vault for a player.
func (v Vault) Open(p *player.Player) {

}

// Clear clears the contents of a vault
func (v Vault) Clear() {
	v.storage = make(map[int]item.Stack)
}

// Contents returns the contents of a vault
func (v Vault) Contents() map[int]item.Stack {
	return v.storage
}
