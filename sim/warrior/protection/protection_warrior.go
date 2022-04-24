package protection

import (
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/warrior"
)

func RegisterProtectionWarrior() {
	core.RegisterAgentFactory(
		proto.Player_ProtectionWarrior{},
		proto.Spec_SpecProtectionWarrior,
		func(character core.Character, options proto.Player) core.Agent {
			return NewProtectionWarrior(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ProtectionWarrior)
			if !ok {
				panic("Invalid spec value for Protection Warrior!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewProtectionWarrior(character core.Character, options proto.Player) *ProtectionWarrior {
	warOptions := options.GetProtectionWarrior()

	war := &ProtectionWarrior{
		Warrior: warrior.NewWarrior(character, *warOptions.Talents, warrior.WarriorInputs{
			ShoutType:            warOptions.Options.Shout,
			PrecastShout:         warOptions.Options.PrecastShout,
			PrecastShoutSapphire: warOptions.Options.PrecastShoutSapphire,
			PrecastShoutT2:       warOptions.Options.PrecastShoutT2,
		}),
		Rotation: *warOptions.Rotation,
		Options:  *warOptions.Options,
	}

	war.EnableRageBar(warOptions.Options.StartingRage, core.TernaryFloat64(war.Talents.EndlessRage, 1.25, 1), func(sim *core.Simulation) {
		if war.GCD.IsReady(sim) {
			war.TryUseCooldowns(sim)
			if war.GCD.IsReady(sim) {
				war.doRotation(sim)
			}
		}
	})
	war.EnableAutoAttacks(war, core.AutoAttackOptions{
		MainHand:       war.WeaponFromMainHand(war.DefaultMeleeCritMultiplier()),
		OffHand:        war.WeaponFromOffHand(war.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		ReplaceMHSwing: func(sim *core.Simulation) *core.Spell {
			if war.Rotation.UseCleave {
				return war.TryCleave(sim)
			} else {
				return war.TryHeroicStrike(sim)
			}
		},
	})

	return war
}

type ProtectionWarrior struct {
	*warrior.Warrior

	Rotation proto.ProtectionWarrior_Rotation
	Options  proto.ProtectionWarrior_Options
}

func (war *ProtectionWarrior) GetWarrior() *warrior.Warrior {
	return war.Warrior
}

func (war *ProtectionWarrior) Reset(sim *core.Simulation) {
	war.Warrior.Reset(sim)
	war.DefensiveStanceAura.Activate(sim)
	war.Stance = warrior.DefensiveStance
}
