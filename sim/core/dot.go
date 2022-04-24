package core

import (
	"time"
)

type TickEffects func(*Simulation, *Spell) func()

type Dot struct {
	Spell *Spell

	// Embed Aura, so we can use IsActive/Refresh/etc directly.
	*Aura

	NumberOfTicks int           // number of ticks over the whole duration
	TickLength    time.Duration // time between each tick

	// If true, tick length will be shortened based on casting speed.
	AffectedByCastSpeed bool

	TickEffects TickEffects

	tickFn     func()
	tickAction *PendingAction
	tickPeriod time.Duration

	// Number of ticks since last call to Apply().
	TickCount int

	lastTickTime time.Duration
}

func (dot *Dot) Apply(sim *Simulation) {
	dot.Cancel(sim)

	dot.TickCount = 0
	if dot.AffectedByCastSpeed {
		castSpeed := dot.Spell.Character.CastSpeed()
		dot.tickPeriod = time.Duration(float64(dot.TickLength) / castSpeed)
		dot.Aura.Duration = dot.tickPeriod*time.Duration(dot.NumberOfTicks) + 1
	}
	dot.Aura.Activate(sim)
}

func (dot *Dot) Cancel(sim *Simulation) {
	if dot.Aura.IsActive() {
		dot.Aura.Deactivate(sim)
	}
}

// Call this after manually changing NumberOfTicks or TickLength.
func (dot *Dot) RecomputeAuraDuration() {
	if dot.AffectedByCastSpeed {
		castSpeed := dot.Spell.Character.CastSpeed()
		dot.tickPeriod = time.Duration(float64(dot.TickLength) / castSpeed)
		dot.Aura.Duration = dot.tickPeriod*time.Duration(dot.NumberOfTicks) + 1
	} else {
		dot.tickPeriod = dot.TickLength
		dot.Aura.Duration = dot.tickPeriod*time.Duration(dot.NumberOfTicks) + 1
	}
}

func (dot *Dot) tick(sim *Simulation) {
	dot.lastTickTime = sim.CurrentTime
	dot.TickCount++
	dot.tickFn()
}

func NewDot(config Dot) *Dot {
	dot := &Dot{}
	*dot = config

	dot.tickPeriod = dot.TickLength
	dot.Aura.Duration = dot.TickLength*time.Duration(dot.NumberOfTicks) + 1

	dot.Aura.OnGain = func(aura *Aura, sim *Simulation) {
		dot.tickFn = dot.TickEffects(sim, dot.Spell)

		pa := &PendingAction{
			NextActionAt: sim.CurrentTime + dot.tickPeriod,
		}
		pa.OnAction = func(sim *Simulation) bool {
			if dot.lastTickTime == sim.CurrentTime {
				return false
			}
			dot.tick(sim)
			pa.NextActionAt = sim.CurrentTime + dot.tickPeriod
			return true
		}

		dot.tickAction = pa
		sim.AddPendingAction(dot.tickAction)
	}
	dot.Aura.OnExpire = func(aura *Aura, sim *Simulation) {
		if dot.tickAction != nil {
			dot.tickAction.Cancel(sim)
			dot.tickAction = nil
		}
	}

	return dot
}

func TickFuncSnapshot(target *Target, baseEffect SpellEffect) TickEffects {
	snapshotEffect := &SpellEffect{}
	return func(sim *Simulation, spell *Spell) func() {
		*snapshotEffect = baseEffect
		snapshotEffect.Target = target
		baseDamage := snapshotEffect.calculateBaseDamage(sim, spell) * snapshotEffect.DamageMultiplier
		snapshotEffect.DamageMultiplier = 1
		snapshotEffect.BaseDamage = BaseDamageConfigFlat(baseDamage)

		effectsFunc := ApplyEffectFuncDirectDamage(*snapshotEffect)
		return func() {
			effectsFunc(sim, target, spell)
		}
	}
}
func TickFuncAOESnapshot(sim *Simulation, baseEffect SpellEffect) TickEffects {
	snapshotEffect := &SpellEffect{}
	return func(sim *Simulation, spell *Spell) func() {
		target := sim.GetPrimaryTarget()
		*snapshotEffect = baseEffect
		snapshotEffect.Target = target
		baseDamage := snapshotEffect.calculateBaseDamage(sim, spell) * snapshotEffect.DamageMultiplier
		snapshotEffect.DamageMultiplier = 1
		snapshotEffect.BaseDamage = BaseDamageConfigFlat(baseDamage)

		effectsFunc := ApplyEffectFuncAOEDamage(sim, *snapshotEffect)
		return func() {
			effectsFunc(sim, target, spell)
		}
	}
}

func TickFuncApplyEffects(effectsFunc ApplySpellEffects) TickEffects {
	return func(sim *Simulation, spell *Spell) func() {
		return func() {
			effectsFunc(sim, sim.GetPrimaryTarget(), spell)
		}
	}
}
