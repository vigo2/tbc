import { RaceBonusType } from './newapi';
import { Spec } from './newapi';
import { Stat } from './newapi';

const shamanRaces = [
    RaceBonusType.draenei,
    RaceBonusType.orc,
    RaceBonusType.tauren,
    RaceBonusType.troll10,
    RaceBonusType.troll30,
];

export const SpecToEligibleRaces: Record<Spec, Array<RaceBonusType>> = {
  [Spec.elemental_shaman]: shamanRaces,
};

export const StatNames: Record<Stat, string> = {
  [Stat.strength]: 'Strength',
  [Stat.agility]: 'Agility',
  [Stat.stamina]: 'Stamina',
  [Stat.intellect]: 'Intellect',
  [Stat.spirit]: 'Spirit',
  [Stat.spell_power]: 'Spell Dmg',
  [Stat.healing_power]: 'Healing Power',
  [Stat.arcane_spell_power]: 'Arcane Dmg',
  [Stat.fire_spell_power]: 'Fire Dmg',
  [Stat.frost_spell_power]: 'Frost Dmg',
  [Stat.holy_spell_power]: 'Holy Dmg',
  [Stat.nature_spell_power]: 'Nature Dmg',
  [Stat.shadow_spell_power]: 'Shadow Dmg',
  [Stat.mp5]: 'MP5',
  [Stat.spell_hit]: 'Spell Hit',
  [Stat.spell_crit]: 'Spell Crit',
  [Stat.spell_haste]: 'Spell Haste',
  [Stat.spell_penetration]: 'Spell Pen',
  [Stat.attack_power]: 'Attack Power',
  [Stat.melee_hit]: 'Hit',
  [Stat.melee_crit]: 'Crit',
  [Stat.melee_haste]: 'Haste',
  [Stat.armor_penetration]: 'Armor Pen',
  [Stat.expertise]: 'Expertise',
  [Stat.mana]: 'Mana',
  [Stat.energy]: 'Energy',
  [Stat.rage]: 'Rage',
  [Stat.armor]: 'Armor',
};