package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var KillCommandCooldownID = core.NewCooldownID()
var KillCommandActionID = core.ActionID{SpellID: 34026, CooldownID: KillCommandCooldownID}

var KillCommandAuraID = core.NewAuraID()

func (hunter *Hunter) applyKillCommand() {
	if hunter.pet == nil {
		return
	}

	hunter.AddPermanentAura(func(sim *core.Simulation) core.Aura {
		return core.Aura{
			ID: KillCommandAuraID,
			OnMeleeAttack: func(sim *core.Simulation, ability *core.ActiveMeleeAbility, hitEffect *core.AbilityHitEffect) {
				if hitEffect.HitType == core.MeleeHitTypeCrit {
					hunter.killCommandEnabled = true
					hunter.TryKillCommand(sim, sim.GetPrimaryTarget())
				}
			},
		}
	})
}

// ActiveMeleeAbility doesn't support cast times, so we wrap it in a SimpleCast.
func (hunter *Hunter) newKillCommandTemplate(sim *core.Simulation) core.SimpleCast {
	template := core.SimpleCast{
		Cast: core.Cast{
			ActionID:     KillCommandActionID,
			Character:    hunter.GetCharacter(),
			BaseManaCost: 75,
			ManaCost:     75,
			Cooldown:     time.Second * 5,
		},
	}

	pa := &core.PendingAction{
		Priority:     core.ActionPriorityRegen,
		NextActionAt: 0,
	}
	pa.OnAction = func(sim *core.Simulation) {
		hunter.TryKillCommand(sim, sim.GetPrimaryTarget())
		hunter.killCommandAction.NextActionAt = 0
	}
	hunter.killCommandAction = pa

	return template
}

func (hp *HunterPet) newKillCommandTemplate(sim *core.Simulation) core.MeleeAbilityTemplate {
	hasBeastLord4Pc := ItemSetBeastLord.CharacterHasSetBonus(&hp.hunterOwner.Character, 4)
	beastLordApplier := hp.hunterOwner.NewTempStatAuraApplier(sim, BeastLord4PcAuraID, core.ActionID{SpellID: 37483}, stats.ArmorPenetration, 600, time.Second*15)

	ama := core.ActiveMeleeAbility{
		MeleeAbility: core.MeleeAbility{
			ActionID:       core.ActionID{SpellID: 34027},
			Character:      &hp.Character,
			SpellSchool:    stats.AttackPower,
			CritMultiplier: 2,
		},
		Effect: core.AbilityHitEffect{
			AbilityEffect: core.AbilityEffect{
				DamageMultiplier:       1,
				StaticDamageMultiplier: 1,
				ThreatMultiplier:       1,
			},
			WeaponInput: core.WeaponDamageInput{
				DamageMultiplier: 1,
				FlatDamageBonus:  127,
			},
		},
		OnMeleeAttack: func(sim *core.Simulation, ability *core.ActiveMeleeAbility, hitEffect *core.AbilityHitEffect) {
			if hasBeastLord4Pc {
				beastLordApplier(sim)
			}
		},
	}

	ama.Effect.BonusCritRating += float64(hp.hunterOwner.Talents.FocusedFire) * 10 * core.MeleeCritRatingPerCritChance

	return core.NewMeleeAbilityTemplate(ama)
}

func (hunter *Hunter) NewKillCommand(sim *core.Simulation, target *core.Target) core.SimpleCast {
	killCommand := hunter.killCommandTemplate

	// Set dynamic fields, i.e. the stuff we couldn't precompute.
	killCommand.OnCastComplete = func(sim *core.Simulation, cast *core.Cast) {
		hunter.killCommandEnabled = false

		kc := &hunter.pet.killCommand
		hunter.pet.killCommandTemplate.Apply(kc)
		kc.Effect.Target = target
		kc.Attack(sim)
	}

	killCommand.Init(sim)
	return killCommand
}

func (hunter *Hunter) TryKillCommand(sim *core.Simulation, target *core.Target) {
	if !hunter.killCommandEnabled || hunter.killCommandBlocked {
		return
	}

	if hunter.CurrentMana() < 75 {
		return
	}

	if !hunter.IsOnCD(KillCommandCooldownID, sim.CurrentTime) {
		kc := hunter.NewKillCommand(sim, target)
		kc.StartCast(sim)
		return
	}

	// Kill Command is on CD, so set up an event to use it when its ready.

	if hunter.killCommandAction.NextActionAt != 0 {
		// An event is already set up from a previous crit.
		return
	}

	hunter.killCommandAction.NextActionAt = hunter.CDReadyAt(KillCommandCooldownID)
	sim.AddPendingAction(hunter.killCommandAction)
}