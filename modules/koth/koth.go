package koth

import (
	"time"
)

// KothEvent is the event in which factions compete together to
// capture the koth. The winning faction gets exclusive prices.
type KothEvent struct {
	// Arena is the Koth Arena in which the Koth Event is taking place
	Arena *KothArena

	// Duration is the total time duration of the Koth Event
	Duration time.Duration

	// Start is unix timestamp of the instant when the Koth started
	Start int64

	// Finish is the unix timestamp of the instant when Koth finished
	Finish int64
}

// Started returns whether the Koth Event has started
func (e *KothEvent) Started() bool {
	return e.Start != 0
}

// Finished returns whether Koth Event has finished
func (e *KothEvent) Finished() bool {
	return e.Finish != 0
}

// TimeLeft returns the duration of time left before the event ends.
func (e *KothEvent) TimeLeft() time.Duration {
	secs := e.Start + int64(e.Duration.Seconds())
	return time.Second * time.Duration(secs)
}
