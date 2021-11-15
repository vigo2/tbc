import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Faction } from '/tbc/core/proto_utils/utils.js';
import { Player } from '/tbc/core/player.js';

import * as Enchants from '/tbc/core/constants/enchants.js';
import * as Gems from '/tbc/core/constants/gems.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://tbc.wowhead.com/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: '510022212503135231351--52003301',
};

export const P1_ALLIANCE_BIS = {
	name: 'P1 Alliance BIS',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Alliance,
	gear: EquipmentSpec.create({
		items: [
			ItemSpec.create({
				id: 29093, // Antlers of Malorne
				enchant: Enchants.GLYPH_OF_POWER,
				gems: [
					Gems.RUNED_LIVING_RUBY,
					Gems.CHAOTIC_SKYFIRE_DIAMOND,
				],
			}),
			ItemSpec.create({
				id: 28762, // Adornment of Stolen Souls
			}),
			ItemSpec.create({
				id: 29095, // Pauldrons of Malorne
				enchant: Enchants.GREATER_INSCRIPTION_OF_DISCIPLINE,
				gems: [
					Gems.GLOWING_NIGHTSEYE,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 28766, // Ruby Drape of the Mysticant
				enchant: Enchants.SUBTLETY,
			}),
			ItemSpec.create({
				id: 21848, // Spellfire Robe
				enchant: Enchants.CHEST_EXCEPTIONAL_STATS,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.GLOWING_NIGHTSEYE,
				],
			}),
			ItemSpec.create({
				id: 24250, // Bracers of Havok
				enchant: Enchants.WRIST_SPELLPOWER,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 21847, // Spellfire Gloves
				enchant: Enchants.GLOVES_SPELLPOWER,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 21846, // Spellfire Belt
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 24262, // Spellstrike Pants
				enchant: Enchants.RUNIC_SPELLTHREAD,
				gems: [
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 28517, // Boots of Foretelling
				enchant: Enchants.BOARS_SPEED,
				gems: [
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 28753, // Ring of Recurrence
				enchant: Enchants.RING_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 29287, // Voilet Signet of the Archmage
				enchant: Enchants.RING_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 29370, // Icon of the Silver Crescent
			}),
			ItemSpec.create({
				id: 27683, // Quagmirran's Eye
			}),
			ItemSpec.create({
				id: 28770, // Nathrezim Mindblade
				enchant: Enchants.SUNFIRE,
			}),
			ItemSpec.create({
				id: 29271, // Talisman of Kalecgos
			}),
			ItemSpec.create({
				id: 27518, // Ivory Idol of the Moongodddess
			}),
		],
	}),
};

export const P1_HORDE_BIS = {
	name: 'P1 Horde BIS',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Horde,
	gear: EquipmentSpec.create({
		items: [
			ItemSpec.create({
				id: 29093, // Antlers of Malorne
				enchant: Enchants.GLYPH_OF_POWER,
				gems: [
					Gems.RUNED_LIVING_RUBY,
					Gems.CHAOTIC_SKYFIRE_DIAMOND,
				],
			}),
			ItemSpec.create({
				id: 28762, // Adornment of Stolen Souls
			}),
			ItemSpec.create({
				id: 29095, // Pauldrons of Malorne
				enchant: Enchants.GREATER_INSCRIPTION_OF_DISCIPLINE,
				gems: [
					Gems.GLOWING_NIGHTSEYE,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 28766, // Ruby Drape of the Mysticant
				enchant: Enchants.SUBTLETY,
			}),
			ItemSpec.create({
				id: 21848, // Spellfire Robe
				enchant: Enchants.CHEST_EXCEPTIONAL_STATS,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.GLOWING_NIGHTSEYE,
				],
			}),
			ItemSpec.create({
				id: 24250, // Bracers of Havok
				enchant: Enchants.WRIST_SPELLPOWER,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 21847, // Spellfire Gloves
				enchant: Enchants.GLOVES_SPELLPOWER,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 21846, // Spellfire Belt
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 24262, // Spellstrike Pants
				enchant: Enchants.RUNIC_SPELLTHREAD,
				gems: [
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 28517, // Boots of Foretelling
				enchant: Enchants.BOARS_SPEED,
				gems: [
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 28753, // Ring of Recurrence
				enchant: Enchants.RING_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 28793, // Band of Crimson Fury
				enchant: Enchants.RING_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 29370, // Icon of the Silver Crescent
			}),
			ItemSpec.create({
				id: 27683, // Quagmirran's Eye
			}),
			ItemSpec.create({
				id: 28770, // Nathrezim Mindblade
				enchant: Enchants.SUNFIRE,
			}),
			ItemSpec.create({
				id: 29271, // Talisman of Kalecgos
			}),
			ItemSpec.create({
				id: 27518, // Ivory Idol of the Moongodddess
			}),
		],
	}),
};

export const P2_ALLIANCE_BIS = {
	name: 'P2 Alliance BIS',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Alliance,
	gear: EquipmentSpec.create({
		items: [
			ItemSpec.create({
				id: 30233, // Nordrassil Headpiece
				enchant: Enchants.GLYPH_OF_POWER,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.CHAOTIC_SKYFIRE_DIAMOND,
				],
			}),
			ItemSpec.create({
				id: 30015, // The Sun King's Talisman
			}),
			ItemSpec.create({
				id: 30235, // Nordrassil Wrath-Mantle
				enchant: Enchants.GREATER_INSCRIPTION_OF_DISCIPLINE,
				gems: [
					Gems.GLOWING_NIGHTSEYE,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 28797, // Brute Cloak of the Ogre-Magi
				enchant: Enchants.SUBTLETY,
			}),
			ItemSpec.create({
				id: 30231, // Nordrassil Chestpiece
				enchant: Enchants.CHEST_EXCEPTIONAL_STATS,
				gems: [
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 29918, // Mindstorm Wristbands
				enchant: Enchants.WRIST_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 21847, // Spellfire Gloves
				enchant: Enchants.GLOVES_SPELLPOWER,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 30038, // Belt of Blasting
				gems: [
					Gems.GLOWING_NIGHTSEYE,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 30234, // Nordrassil Wrath-Kilt
				enchant: Enchants.RUNIC_SPELLTHREAD,
				gems: [
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 30037, // Boots of Blasting
				enchant: Enchants.BOARS_SPEED,
			}),
			ItemSpec.create({
				id: 28753, // Ring of Recurrence
				enchant: Enchants.RING_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 29302, // Band of Eternity
				enchant: Enchants.RING_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 29370, // Icon of the Silver Crescent
			}),
			ItemSpec.create({
				id: 27683, // Quagmirran's Eye
			}),
			ItemSpec.create({
				id: 29988, // The Nexus Key
				enchant: Enchants.SUNFIRE,
			}),
			ItemSpec.create({
				id: 32387, // Idol of the Raven Goddess
			}),
		],
	}),
};

export const P2_HORDE_BIS = {
	name: 'P2 Horde BIS',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<any>) => player.getFaction() == Faction.Horde,
	gear: EquipmentSpec.create({
		items: [
			ItemSpec.create({
				id: 30233, // Nordrassil Headpiece
				enchant: Enchants.GLYPH_OF_POWER,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.CHAOTIC_SKYFIRE_DIAMOND,
				],
			}),
			ItemSpec.create({
				id: 30015, // The Sun King's Talisman
			}),
			ItemSpec.create({
				id: 30235, // Nordrassil Wrath-Mantle
				enchant: Enchants.GREATER_INSCRIPTION_OF_DISCIPLINE,
				gems: [
					Gems.GLOWING_NIGHTSEYE,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 28797, // Brute Cloak of the Ogre-Magi
				enchant: Enchants.SUBTLETY,
			}),
			ItemSpec.create({
				id: 30107, // Vestments of the Sea-Witch
				enchant: Enchants.CHEST_EXCEPTIONAL_STATS,
				gems: [
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 29918, // Mindstorm Wristbands
				enchant: Enchants.WRIST_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 30232, // Nordrassil Gauntlets
				enchant: Enchants.GLOVES_SPELLPOWER,
				gems: [
					Gems.POTENT_NOBLE_TOPAZ,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 30038, // Belt of Blasting
				gems: [
					Gems.GLOWING_NIGHTSEYE,
					Gems.POTENT_NOBLE_TOPAZ,
				],
			}),
			ItemSpec.create({
				id: 30234, // Nordrassil Wrath-Kilt
				enchant: Enchants.RUNIC_SPELLTHREAD,
				gems: [
					Gems.RUNED_LIVING_RUBY,
				],
			}),
			ItemSpec.create({
				id: 30037, // Boots of Blasting
				enchant: Enchants.BOARS_SPEED,
			}),
			ItemSpec.create({
				id: 28753, // Ring of Recurrence
				enchant: Enchants.RING_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 29302, // Band of Eternity
				enchant: Enchants.RING_SPELLPOWER,
			}),
			ItemSpec.create({
				id: 29370, // Icon of the Silver Crescent
			}),
			ItemSpec.create({
				id: 27683, // Quagmirran's Eye
			}),
			ItemSpec.create({
				id: 29988, // The Nexus Key
				enchant: Enchants.SUNFIRE,
			}),
			ItemSpec.create({
				id: 32387, // Idol of the Raven Goddess
			}),
		],
	}),
};