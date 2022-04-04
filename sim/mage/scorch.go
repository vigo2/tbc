package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

const SpellIDScorch int32 = 27074

func (mage *Mage) registerScorchSpell(sim *core.Simulation) {
	spell := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				ActionID:    core.ActionID{SpellID: SpellIDScorch},
				Character:   &mage.Character,
				SpellSchool: core.SpellSchoolFire,
				SpellExtras: SpellFlagMage,
				BaseCost: core.ResourceCost{
					Type:  stats.Mana,
					Value: 180,
				},
				Cost: core.ResourceCost{
					Type:  stats.Mana,
					Value: 180,
				},
				CastTime: time.Millisecond * 1500,
				GCD:      core.GCDDefault,
			},
		},
		Effect: core.SpellEffect{
			OutcomeRollCategory: core.OutcomeRollCategoryMagic,
			CritRollCategory:    core.CritRollCategoryMagical,
			CritMultiplier:      mage.SpellCritMultiplier(1, 0.25*float64(mage.Talents.SpellPower)),
			DamageMultiplier:    mage.spellDamageMultiplier,
			ThreatMultiplier:    1 - 0.05*float64(mage.Talents.BurningSoul),
			BaseDamage:          core.BaseDamageConfigMagic(305, 361, 1.5/3.5),
		},
	}

	spell.Cost.Value -= spell.BaseCost.Value * float64(mage.Talents.Pyromaniac) * 0.01
	spell.Cost.Value *= 1 - float64(mage.Talents.ElementalPrecision)*0.01
	spell.Effect.BonusSpellHitRating += float64(mage.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance
	spell.Effect.BonusSpellCritRating += float64(mage.Talents.Incineration) * 2 * core.SpellCritRatingPerCritChance
	spell.Effect.BonusSpellCritRating += float64(mage.Talents.CriticalMass) * 2 * core.SpellCritRatingPerCritChance
	spell.Effect.BonusSpellCritRating += float64(mage.Talents.Pyromaniac) * 1 * core.SpellCritRatingPerCritChance
	spell.Effect.DamageMultiplier *= 1 + 0.02*float64(mage.Talents.FirePower)

	if mage.Talents.ImprovedScorch > 0 {
		mage.ScorchAura = sim.GetPrimaryTarget().GetAura(core.ImprovedScorchAuraLabel)
		if mage.ScorchAura == nil {
			mage.ScorchAura = core.ImprovedScorchAura(sim.GetPrimaryTarget(), 0)
		}

		procChance := float64(mage.Talents.ImprovedScorch) / 3.0
		spell.Effect.OnSpellHit = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}

			if procChance != 1.0 || sim.RandomFloat("Improved Scorch") > procChance {
				return
			}

			mage.ScorchAura.Activate(sim)
			mage.ScorchAura.AddStack(sim)
		}
	}

	mage.Scorch = mage.RegisterSpell(core.SpellConfig{
		Template:   spell,
		ModifyCast: core.ModifyCastAssignTarget,
	})
}
