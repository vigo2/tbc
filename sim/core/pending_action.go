package core

import (
	"time"
)

type ActionPriority int32

const (
	ActionPriorityLow ActionPriority = -1
	ActionPriorityGCD ActionPriority = 0

	// Higher than GCD because regen can cause GCD actions (if we were waiting
	// for mana).
	ActionPriorityRegen ActionPriority = 1

	// Autos can cause regen (JoW, rage, energy procs, etc) so they should be
	// higher prio so that we never go backwards in the priority order.
	ActionPriorityAuto ActionPriority = 2

	// DOTs need to be higher than anything else so that dots can properly expire before we take other actions.
	ActionPriorityDOT ActionPriority = 3
)

type PendingAction struct {
	NextActionAt time.Duration
	OnAction     func(*Simulation) bool

	Priority  ActionPriority
	cancelled bool
}

func (pa *PendingAction) Cancel(_ *Simulation) {
	pa.cancelled = true
}
