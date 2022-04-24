package core

import (
	"time"
)

type PeriodicActionOptions struct {
	// How often the action should be performed.
	Period time.Duration

	// Number of times to perform the action before stopping.
	// 0 indicates a permanent periodic action.
	NumTicks int

	OnAction func(*Simulation)
}

func NewPeriodicAction(sim *Simulation, options PeriodicActionOptions) *PendingAction {
	pa := &PendingAction{
		NextActionAt: sim.CurrentTime + options.Period,
	}

	tickIndex := 0

	pa.OnAction = func(sim *Simulation) bool {
		options.OnAction(sim)
		tickIndex++

		if options.NumTicks != 0 && tickIndex >= options.NumTicks {
			return false
		}

		// Refresh action.
		pa.NextActionAt = sim.CurrentTime + options.Period
		return true
	}

	return pa
}

// Convenience for immediately creating and starting a periodic action.
func StartPeriodicAction(sim *Simulation, options PeriodicActionOptions) {
	pa := NewPeriodicAction(sim, options)
	sim.AddPendingAction(pa)
}
