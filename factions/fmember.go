package factions

type FMember struct {
	// Name is the name of the Faction Member
	Name string `json:"name"`

	// Xuid is the xuid of the Faction Member
	Xuid string `json:"xuid"`

	// Rank is the rank of the Faction member
	Rank int `json:"Rank"`
}

const (
	Recruit = iota
	Assistant
	Officer
	Manager
	CoLeader
	Leader
)

// Compare returns whether the faction member is above a Faction member in hierarchy
func (m *FMember) Compare(t *FMember) bool {
	return m.Rank > t.Rank
}

// Successor returns the next faction rank
func (m *FMember) Successor() int {
	if m.Rank == 5 {
		return 5
	}

	return m.Rank + 1
}

// Predecessor returns the previous faction rank
func (m *FMember) Predecessor() int {
	if m.Rank == 0 {
		return 0
	}

	return m.Rank - 1
}
