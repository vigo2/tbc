package druid

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Starfire spell IDs
const SpellIDWrath int32 = 26985

func (druid *Druid) newWrathTemplate(sim *core.Simulation) core.SingleTargetDirectDamageSpellTemplate {
	baseCast := core.Cast{
		Name:           "Wrath",
		CritMultiplier: 1.5,
		SpellSchool:    stats.NatureSpellPower,
		Character:      &druid.Character,
		BaseManaCost:   255,
		ManaCost:       255,
		CastTime:       time.Millisecond * 2000,
		ActionID: core.ActionID{
			SpellID: SpellIDWrath,
		},
	}

	effect := core.DirectDamageSpellEffect{
		SpellEffect: core.SpellEffect{
			DamageMultiplier: 1,
		},
		DirectDamageSpellInput: core.DirectDamageSpellInput{
			MinBaseDamage:    383,
			MaxBaseDamage:    432,
			SpellCoefficient: 0.571,
		},
	}

	// TODO: Applies to both starfire and moonfire
	baseCast.CastTime -= time.Millisecond * 100 * time.Duration(druid.Talents.StarlightWrath)
	effect.SpellEffect.BonusSpellCritRating += float64(druid.Talents.FocusedStarlight) * 2 * core.SpellCritRatingPerCritChance // 2% crit per point

	// TODO: applies to starfire, wrath and moonfire

	// Convert to percent, multiply by percent increase, convert back to multiplier by adding 1
	baseCast.CritMultiplier = (baseCast.CritMultiplier-1)*(1+float64(druid.Talents.Vengeance)*0.2) + 1
	baseCast.ManaCost -= baseCast.BaseManaCost * float64(druid.Talents.Moonglow) * 0.03
	effect.SpellEffect.DamageMultiplier *= 1 + 0.02*float64(druid.Talents.Moonfury)
	effect.SpellEffect.BonusSpellHitRating += float64(druid.Talents.BalanceOfPower) * 2 * core.SpellHitRatingPerHitChance

	effect.OnSpellHit = druid.applyOnHitTalents
	spCast := &core.SpellCast{
		Cast: baseCast,
	}

	return core.NewSingleTargetDirectDamageSpellTemplate(core.SingleTargetDirectDamageSpell{
		SpellCast: *spCast,
		Effect:    effect,
	})
}

func (druid *Druid) NewWrath(sim *core.Simulation, target *core.Target) *core.SingleTargetDirectDamageSpell {
	// Initialize cast from precomputed template.
	sf := &druid.wrathSpell

	druid.wrathCastTemplate.Apply(sf)

	// Modifies the cast time.
	druid.applyNaturesGrace(&sf.SpellCast)

	// Set dynamic fields, i.e. the stuff we couldn't precompute.
	sf.Effect.Target = target
	sf.Init(sim)

	return sf
}