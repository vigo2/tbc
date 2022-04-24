package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

const SpellIDSiphonLife6 int32 = 30911

var SiphonLife6ActionID = core.ActionID{SpellID: SpellIDSiphonLife6}

func (warlock *Warlock) registerSiphonLifeSpell(sim *core.Simulation) {
	baseCost := 370.0
	effect := core.SpellEffect{
		OutcomeApplier: core.OutcomeFuncMagicHit(),
		OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() {
				warlock.SiphonLifeDot.Apply(sim)
			}
		},
	}

	warlock.SiphonLife = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     SiphonLife6ActionID,
		SpellSchool:  core.SpellSchoolShadow,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	target := sim.GetPrimaryTarget()
	spellCoefficient := 0.1

	warlock.SiphonLifeDot = core.NewDot(core.Dot{
		Spell: warlock.SiphonLife,
		Aura: target.RegisterAura(core.Aura{
			Label:    "SiphonLife-" + strconv.Itoa(int(warlock.Index)),
			ActionID: SiphonLife6ActionID,
		}),
		NumberOfTicks: 10,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier: 1 * (1 + 0.02*float64(warlock.Talents.ShadowMastery)) * (1 + 0.01*float64(warlock.Talents.Contagion)),
			ThreatMultiplier: 1 - 0.05*float64(warlock.Talents.ImprovedDrainSoul),
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(63, spellCoefficient),
			OutcomeApplier:   core.OutcomeFuncTick(),
			IsPeriodic:       true,
		}),
	})
}
