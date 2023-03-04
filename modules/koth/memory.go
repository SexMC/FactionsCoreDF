package koth

// Koth is the current on-going Koth Event. It may be nil if there's currently
// no Koth going on
var Koth *KothEvent

// Unix timestamp when the Next Koth Event will start.
var NextKoth int64

// Status returns whether a Koth event is either scheduled or
// running
func Status() bool {
	return Koth != nil
}
