package chat

type Channel interface {
	// ChannelType returns the type of the channel
	ChannelType() int
}

const (
	Global = iota
	Truces
	Allies
	Faction
	Moderator
)

var channelNames = map[int]string{
	Global:    "Global",
	Truces:    "Truces",
	Allies:    "Allies",
	Faction:   "Faction",
	Moderator: "Moderator",
}

var channelIDs = map[int]string{
	Global:    "global",
	Truces:    "truces",
	Allies:    "allies",
	Faction:   "faction",
	Moderator: "moderator",
}

// GlobalChannel is the default chat channel, all messages sent to this
// channel are visible to all the players
type GlobalChannel struct {
	Type int
}

// ChannelType ...
func (GlobalChannel) ChannelType() int {
	return Global
}

// TrucesChannel is the chat channel where all the messages are visible to the truces
// of a faction
type TrucesChannel struct {
	Type int
}

// ChannelType ...
func (TrucesChannel) ChannelType() int {
	return Truces
}

// AlliesChannel is the chat channel where all the messages are visible to the allies
// of a faction
type AlliesChannel struct {
	Type int
}

// ChannelType ...
func (AlliesChannel) ChannelType() int {
	return Allies
}

// FactionChannel is the chat channel where all the messages are visible to the members
// of a faction
type FactionChannel struct {
	Type int
}

// ChannelType ...
func (FactionChannel) ChannelType() int {
	return Faction
}

// ModeratorChannel is the chat channel where all the messages are visible to the managers+
// of a faction
type ModeratorChannel struct {
	Type int
}

// ChannelType ...
func (ModeratorChannel) ChannelType() int {
	return Moderator
}

// ChannelName returns the channel name of a channel
func ChannelName(channel Channel) string {
	return channelNames[channel.ChannelType()]
}

// ChannelID returns the channel id of a channel
func ChannelID(channel Channel) string {
	return channelIDs[channel.ChannelType()]
}
