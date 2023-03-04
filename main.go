package main

import (
	"github.com/inceptionmc/factions/commands"
	"github.com/inceptionmc/factions/memory"
	"github.com/inceptionmc/factions/postgres"
	"github.com/linuxtf/dragonfly/libraries/console"
	"github.com/linuxtf/dragonfly/libraries/srv"
	"github.com/linuxtf/dragonfly/server/cmd"
	"github.com/linuxtf/dragonfly/server/player"
)

func main() {
	// Register Commands
	registerCommands()

	// Start the Server
	srv.Start()

	// Initiate Database
	postgres.Init()

	// Load Faction Claims
	c := memory.LoadClaims()
	console.Log.Printf("[Factions] Loaded %d claims!", c)

	// Start Listening for Connections
	for srv.Srv.Accept(onPlayerJoin) {
	}
}

func onPlayerJoin(p *player.Player) {
	srv.PlayerJoined(p)

	// Register Handlers
	p.Handle(Handler(p))

	// Load Player Data
	memory.LoadFPlayer(p)
}

// registerCommands registers all the commands
func registerCommands() {
	cmd.Register(cmd.New("faction", "Faction commands", []string{"f"}, commands.FCreateCmd{}, commands.FDescCmd{}, commands.FDisbandCmd{}, commands.FInviteCmd{}, commands.FInfoCmd{}, commands.FPromoteCmd{}, commands.FDemoteCmd{}, commands.FChatCmd{}, commands.FLeaveCmd{}, commands.FInvitesCmd{}, commands.FKickCmd{}, commands.FTruceCmd{}, commands.FAllyCmd{}, commands.FNeutralCmd{}, commands.FEnemyCmd{}, commands.FHomeCmd{}, commands.FWarpCmd{}, commands.FSetHomeCmd{}, commands.FSetWarpCmd{}, commands.FAlertCmd{}, commands.FClaimCmd{}, commands.FAutoClaimCmd{}, commands.FUnclaimCmd{}, commands.FUnclaimAllCmd{}, commands.FMapCmd{}, commands.FRecruitCmd{}, commands.FAssistantCmd{}, commands.FOfficerCmd{}, commands.FManagerCmd{}, commands.FColeaderCmd{}))
	cmd.Register(cmd.New("test", "Test command", []string{}, commands.TestCmd{}))
}
