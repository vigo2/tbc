import { Debuffs } from '/tbc/core/proto/common.js';
import { MobType } from '/tbc/core/proto/common.js';
import { Target as TargetProto } from '/tbc/core/proto/common.js';

import { Listener } from './typed_event.js';
import { Sim } from './sim.js';
import { TypedEvent } from './typed_event.js';
import { sum } from './utils.js';
import { wait } from './utils.js';

export interface TargetConfig {
  defaults: {
		armor: number,
		mobType: MobType,
		debuffs: Debuffs,
  },
}

// Manages all the settings for a single Target.
export class Target {
  readonly armorChangeEmitter = new TypedEvent<void>();
  readonly mobTypeChangeEmitter = new TypedEvent<void>();
  readonly debuffsChangeEmitter = new TypedEvent<void>();

  // Emits when any of the above emitters emit.
  readonly changeEmitter = new TypedEvent<void>();

  // Current values
	private armor: number;
	private mobType: MobType;
  private debuffs: Debuffs;

	private readonly sim: Sim;

  constructor(config: TargetConfig, sim: Sim) {
		this.sim = sim;

    this.armor = config.defaults.armor;
    this.mobType = config.defaults.mobType;
    this.debuffs = config.defaults.debuffs;

    [
      this.armorChangeEmitter,
      this.mobTypeChangeEmitter,
      this.debuffsChangeEmitter,
    ].forEach(emitter => emitter.on(() => this.changeEmitter.emit()));
  }

  getArmor(): number {
    return this.armor;
  }

  setArmor(newArmor: number) {
    if (newArmor == this.armor)
      return;

		this.armor = newArmor;
    this.armorChangeEmitter.emit();
  }

  getMobType(): MobType {
    return this.mobType;
  }

  setMobType(newMobType: MobType) {
    if (newMobType == this.mobType)
      return;

		this.mobType = newMobType;
    this.mobTypeChangeEmitter.emit();
  }

  getDebuffs(): Debuffs {
    // Make a defensive copy
    return Debuffs.clone(this.debuffs);
  }

  setDebuffs(newDebuffs: Debuffs) {
    if (Debuffs.equals(this.debuffs, newDebuffs))
      return;

    // Make a defensive copy
    this.debuffs = Debuffs.clone(newDebuffs);
    this.debuffsChangeEmitter.emit();
  }

	toProto(): TargetProto {
		return TargetProto.create({
			armor: this.armor,
			mobType: this.mobType,
			debuffs: this.debuffs,
		});
	}

  // Returns JSON representing all the current values.
  toJson(): Object {
    return {
      'armor': this.armor,
      'mobType': this.mobType,
      'debuffs': Debuffs.toJson(this.debuffs),
    };
  }

  // Set all the current values, assumes obj is the same type returned by toJson().
  fromJson(obj: any) {
		const parsedArmor = parseInt(obj['armor']);
		if (!isNaN(parsedArmor) && parsedArmor != 0) {
			this.setArmor(parsedArmor);
		}

		const parsedMobType = parseInt(obj['mobType']);
		if (!isNaN(parsedMobType) && parsedMobType != 0) {
			this.setMobType(parsedMobType);
		}

		try {
			this.setDebuffs(Debuffs.fromJson(obj['debuffs']));
		} catch (e) {
			console.warn('Failed to parse debuffs: ' + e);
		}
  }
}