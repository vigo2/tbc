import { Consumes } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { Player as PlayerProto } from '/tbc/core/proto/api.js';
import { PlayerOptions as PlayerOptionsProto } from '/tbc/core/proto/api.js';
import { ComputeStatsRequest, ComputeStatsResult } from '/tbc/core/proto/api.js';
import { Gear } from '/tbc/core/proto_utils/gear.js';
import { Stats } from '/tbc/core/proto_utils/stats.js';
import { canEquipItem, getEligibleItemSlots, getMetaGemEffectEP, gemMatchesSocket, raceToFaction, specToClass, specToEligibleRaces, specTypeFunctions, withSpecProto, } from '/tbc/core/proto_utils/utils.js';
import { TypedEvent } from './typed_event.js';
import { sum } from './utils.js';
import { wait } from './utils.js';
// Manages all the gear / consumes / other settings for a single Player.
export class Player {
    constructor(spec, sim) {
        this.consumes = Consumes.create();
        this.customStats = new Stats();
        this.gear = new Gear({});
        this.talentsString = '';
        this.epWeights = new Stats();
        this.consumesChangeEmitter = new TypedEvent();
        this.customStatsChangeEmitter = new TypedEvent();
        this.gearChangeEmitter = new TypedEvent();
        this.raceChangeEmitter = new TypedEvent();
        this.rotationChangeEmitter = new TypedEvent();
        this.talentsChangeEmitter = new TypedEvent();
        // Talents dont have all fields so we need this.
        this.talentsStringChangeEmitter = new TypedEvent();
        this.specOptionsChangeEmitter = new TypedEvent();
        this.currentStatsEmitter = new TypedEvent();
        // Emits when any of the above emitters emit.
        this.changeEmitter = new TypedEvent();
        this.sim = sim;
        this.spec = spec;
        this.race = specToEligibleRaces[this.spec][0];
        this.specTypeFunctions = specTypeFunctions[this.spec];
        this.rotation = this.specTypeFunctions.rotationCreate();
        this.talents = this.specTypeFunctions.talentsCreate();
        this.specOptions = this.specTypeFunctions.optionsCreate();
        [
            this.consumesChangeEmitter,
            this.customStatsChangeEmitter,
            this.gearChangeEmitter,
            this.raceChangeEmitter,
            this.rotationChangeEmitter,
            this.talentsChangeEmitter,
            this.talentsStringChangeEmitter,
            this.specOptionsChangeEmitter,
        ].forEach(emitter => emitter.on(() => this.changeEmitter.emit()));
        this.currentStats = ComputeStatsResult.create();
        this.sim.changeEmitter.on(() => {
            this.updateCharacterStats();
        });
        this.changeEmitter.on(() => {
            this.updateCharacterStats();
        });
    }
    // Returns all items that this player can wear in the given slot.
    getItems(slot) {
        return this.sim.getItems(slot).filter(item => canEquipItem(item, this.spec));
    }
    // Returns all enchants that this player can wear in the given slot.
    getEnchants(slot) {
        return this.sim.getEnchants(slot);
    }
    // Returns all gems that this player can wear of the given color.
    getGems(socketColor) {
        return this.sim.getGems(socketColor);
    }
    getEpWeights() {
        return this.epWeights;
    }
    setEpWeights(newEpWeights) {
        this.epWeights = newEpWeights;
    }
    async statWeights(request) {
        const result = await this.sim.statWeights(request);
        this.epWeights = new Stats(result.epValues);
        return result;
    }
    // This should be invoked internally whenever stats might have changed.
    async updateCharacterStats() {
        // Sometimes a ui change triggers other changes, so waiting a bit makes sure
        // we get all of them.
        await wait(10);
        const computeStatsResult = await this.sim.computeStats(ComputeStatsRequest.create({
            player: this.toProto(),
            raidBuffs: this.sim.getRaidBuffs(),
            partyBuffs: this.sim.getPartyBuffs(),
            individualBuffs: this.sim.getIndividualBuffs(),
        }));
        this.currentStats = computeStatsResult;
        this.currentStatsEmitter.emit();
    }
    getCurrentStats() {
        return ComputeStatsResult.clone(this.currentStats);
    }
    getRace() {
        return this.race;
    }
    setRace(newRace) {
        if (newRace != this.race) {
            this.race = newRace;
            this.raceChangeEmitter.emit();
        }
    }
    getFaction() {
        return raceToFaction[this.getRace()];
    }
    getConsumes() {
        // Make a defensive copy
        return Consumes.clone(this.consumes);
    }
    setConsumes(newConsumes) {
        if (Consumes.equals(this.consumes, newConsumes))
            return;
        // Make a defensive copy
        this.consumes = Consumes.clone(newConsumes);
        this.consumesChangeEmitter.emit();
    }
    equipItem(slot, newItem) {
        const newGear = this.gear.withEquippedItem(slot, newItem);
        if (newGear.equals(this.gear))
            return;
        this.gear = newGear;
        this.gearChangeEmitter.emit();
    }
    getEquippedItem(slot) {
        return this.gear.getEquippedItem(slot);
    }
    getGear() {
        return this.gear;
    }
    setGear(newGear) {
        if (newGear.equals(this.gear))
            return;
        this.gear = newGear;
        this.gearChangeEmitter.emit();
    }
    getCustomStats() {
        return this.customStats;
    }
    setCustomStats(newCustomStats) {
        if (newCustomStats.equals(this.customStats))
            return;
        this.customStats = newCustomStats;
        this.customStatsChangeEmitter.emit();
    }
    getRotation() {
        return this.specTypeFunctions.rotationCopy(this.rotation);
    }
    setRotation(newRotation) {
        if (this.specTypeFunctions.rotationEquals(newRotation, this.rotation))
            return;
        this.rotation = this.specTypeFunctions.rotationCopy(newRotation);
        this.rotationChangeEmitter.emit();
    }
    getTalents() {
        return this.specTypeFunctions.talentsCopy(this.talents);
    }
    setTalents(newTalents) {
        if (this.specTypeFunctions.talentsEquals(newTalents, this.talents))
            return;
        this.talents = this.specTypeFunctions.talentsCopy(newTalents);
        this.talentsChangeEmitter.emit();
    }
    getTalentsString() {
        return this.talentsString;
    }
    setTalentsString(newTalentsString) {
        if (newTalentsString == this.talentsString)
            return;
        this.talentsString = newTalentsString;
        this.talentsStringChangeEmitter.emit();
    }
    getSpecOptions() {
        return this.specTypeFunctions.optionsCopy(this.specOptions);
    }
    setSpecOptions(newSpecOptions) {
        if (this.specTypeFunctions.optionsEquals(newSpecOptions, this.specOptions))
            return;
        this.specOptions = this.specTypeFunctions.optionsCopy(newSpecOptions);
        this.specOptionsChangeEmitter.emit();
    }
    computeGemEP(gem) {
        const epFromStats = new Stats(gem.stats).computeEP(this.epWeights);
        const epFromEffect = getMetaGemEffectEP(this.spec, gem, new Stats(this.currentStats.finalStats));
        return epFromStats + epFromEffect;
    }
    computeEnchantEP(enchant) {
        return new Stats(enchant.stats).computeEP(this.epWeights);
    }
    computeItemEP(item) {
        if (item == null)
            return 0;
        let ep = new Stats(item.stats).computeEP(this.epWeights);
        const slot = getEligibleItemSlots(item)[0];
        const enchants = this.sim.getEnchants(slot);
        if (enchants.length > 0) {
            ep += Math.max(...enchants.map(enchant => this.computeEnchantEP(enchant)));
        }
        // Compare whether its better to match sockets + get socket bonus, or just use best gems.
        const bestGemEPNotMatchingSockets = sum(item.gemSockets.map(socketColor => {
            const gems = this.sim.getGems(socketColor).filter(gem => !gem.unique && gem.phase <= this.sim.getPhase());
            if (gems.length > 0) {
                return Math.max(...gems.map(gem => this.computeGemEP(gem)));
            }
            else {
                return 0;
            }
        }));
        const bestGemEPMatchingSockets = sum(item.gemSockets.map(socketColor => {
            const gems = this.sim.getGems(socketColor).filter(gem => !gem.unique && gem.phase <= this.sim.getPhase() && gemMatchesSocket(gem, socketColor));
            if (gems.length > 0) {
                return Math.max(...gems.map(gem => this.computeGemEP(gem)));
            }
            else {
                return 0;
            }
        })) + new Stats(item.socketBonus).computeEP(this.epWeights);
        ep += Math.max(bestGemEPMatchingSockets, bestGemEPNotMatchingSockets);
        return ep;
    }
    setWowheadData(equippedItem, elem) {
        let parts = [];
        if (equippedItem.gems.length > 0) {
            parts.push('gems=' + equippedItem.gems.map(gem => gem ? gem.id : 0).join(':'));
        }
        if (equippedItem.enchant != null) {
            parts.push('ench=' + equippedItem.enchant.effectId);
        }
        parts.push('pcs=' + this.gear.asArray().filter(ei => ei != null).map(ei => ei.item.id).join(':'));
        elem.setAttribute('data-wowhead', parts.join('&'));
    }
    toProto() {
        return PlayerProto.create({
            customStats: this.getCustomStats().asArray(),
            equipment: this.getGear().asSpec(),
            options: withSpecProto(PlayerOptionsProto.create({
                race: this.getRace(),
                class: specToClass[this.spec],
                consumes: this.getConsumes(),
            }), this.getRotation(), this.getTalents(), this.getSpecOptions()),
        });
    }
    // TODO: Remove to/from json functions and use proto versions instead. This will require
    // some way to store all talents in the proto.
    // Returns JSON representing all the current values.
    toJson() {
        return {
            'consumes': Consumes.toJson(this.consumes),
            'customStats': this.customStats.toJson(),
            'gear': EquipmentSpec.toJson(this.gear.asSpec()),
            'race': this.race,
            'rotation': this.specTypeFunctions.rotationToJson(this.rotation),
            'talents': this.talentsString,
            'specOptions': this.specTypeFunctions.optionsToJson(this.specOptions),
        };
    }
    // Set all the current values, assumes obj is the same type returned by toJson().
    fromJson(obj) {
        try {
            this.setConsumes(Consumes.fromJson(obj['consumes']));
        }
        catch (e) {
            console.warn('Failed to parse consumes: ' + e);
        }
        try {
            this.setCustomStats(Stats.fromJson(obj['customStats']));
        }
        catch (e) {
            console.warn('Failed to parse custom stats: ' + e);
        }
        try {
            this.setGear(this.sim.lookupEquipmentSpec(EquipmentSpec.fromJson(obj['gear'])));
        }
        catch (e) {
            console.warn('Failed to parse gear: ' + e);
        }
        try {
            this.setRace(obj['race']);
        }
        catch (e) {
            console.warn('Failed to parse race: ' + e);
        }
        try {
            this.setRotation(this.specTypeFunctions.rotationFromJson(obj['rotation']));
        }
        catch (e) {
            console.warn('Failed to parse rotation: ' + e);
        }
        try {
            this.setTalentsString(obj['talents']);
        }
        catch (e) {
            console.warn('Failed to parse talents: ' + e);
        }
        try {
            this.setSpecOptions(this.specTypeFunctions.optionsFromJson(obj['specOptions']));
        }
        catch (e) {
            console.warn('Failed to parse spec options: ' + e);
        }
    }
}