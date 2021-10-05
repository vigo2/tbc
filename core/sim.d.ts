import { EquippedItem } from './api/equipped_item.js';
import { Gear } from './api/gear.js';
import { Buffs } from './proto/common.js';
import { Consumes } from './proto/common.js';
import { Enchant } from './proto/common.js';
import { Encounter } from './proto/common.js';
import { EquipmentSpec } from './proto/common.js';
import { Gem } from './proto/common.js';
import { GemColor } from './proto/common.js';
import { ItemSlot } from './proto/common.js';
import { ItemSpec } from './proto/common.js';
import { Item } from './proto/common.js';
import { Race } from './proto/common.js';
import { Spec } from './proto/common.js';
import { Stat } from './proto/common.js';
import { Stats } from './api/stats.js';
import { SpecAgent } from './api/utils.js';
import { SpecTalents } from './api/utils.js';
import { SpecTypeFunctions } from './api/utils.js';
import { SpecOptions } from './api/utils.js';
import { ComputeStatsResult } from './proto/api.js';
import { IndividualSimRequest } from './proto/api.js';
import { StatWeightsRequest, StatWeightsResult } from './proto/api.js';
import { TypedEvent } from './typed_event.js';
import { WorkerPool } from './worker_pool.js';
export interface SimConfig<SpecType extends Spec> {
    spec: Spec;
    epStats: Array<Stat>;
    epReferenceStat: Stat;
    defaults: {
        phase: number;
        epWeights: Stats;
        encounter: Encounter;
        buffs: Buffs;
        consumes: Consumes;
        agent: SpecAgent<SpecType>;
        talents: string;
        specOptions: SpecOptions<SpecType>;
    };
    metaGemEffectEP?: ((gem: Gem, sim: Sim<SpecType>) => number);
}
export declare class Sim<SpecType extends Spec> extends WorkerPool {
    readonly spec: Spec;
    readonly phaseChangeEmitter: TypedEvent<void>;
    readonly buffsChangeEmitter: TypedEvent<void>;
    readonly consumesChangeEmitter: TypedEvent<void>;
    readonly customStatsChangeEmitter: TypedEvent<void>;
    readonly encounterChangeEmitter: TypedEvent<void>;
    readonly gearChangeEmitter: TypedEvent<void>;
    readonly raceChangeEmitter: TypedEvent<void>;
    readonly agentChangeEmitter: TypedEvent<void>;
    readonly talentsChangeEmitter: TypedEvent<void>;
    readonly talentsStringChangeEmitter: TypedEvent<void>;
    readonly specOptionsChangeEmitter: TypedEvent<void>;
    readonly changeEmitter: TypedEvent<void>;
    readonly gearListEmitter: TypedEvent<void>;
    readonly characterStatsEmitter: TypedEvent<void>;
    private _currentStats;
    private _items;
    private _enchants;
    private _gems;
    private _phase;
    private _buffs;
    private _consumes;
    private _customStats;
    private _gear;
    private _encounter;
    private _race;
    private _agent;
    private _talents;
    private _talentsString;
    private _specOptions;
    private _epWeights;
    readonly specTypeFunctions: SpecTypeFunctions<SpecType>;
    private readonly _metaGemEffectEP;
    private _init;
    constructor(config: SimConfig<SpecType>);
    init(): Promise<void>;
    statWeights(request: StatWeightsRequest): Promise<StatWeightsResult>;
    private updateCharacterStats;
    getCurrentStats(): ComputeStatsResult;
    getItems(slot: ItemSlot | undefined): Array<Item>;
    getEnchants(slot: ItemSlot | undefined): Array<Enchant>;
    getGems(socketColor: GemColor | undefined): Array<Gem>;
    getMatchingGems(socketColor: GemColor): Array<Gem>;
    getPhase(): number;
    setPhase(newPhase: number): void;
    getRace(): Race;
    setRace(newRace: Race): void;
    getBuffs(): Buffs;
    setBuffs(newBuffs: Buffs): void;
    getConsumes(): Consumes;
    setConsumes(newConsumes: Consumes): void;
    getEncounter(): Encounter;
    setEncounter(newEncounter: Encounter): void;
    equipItem(slot: ItemSlot, newItem: EquippedItem | null): void;
    getEquippedItem(slot: ItemSlot): EquippedItem | null;
    getGear(): Gear;
    setGear(newGear: Gear): void;
    getCustomStats(): Stats;
    setCustomStats(newCustomStats: Stats): void;
    getAgent(): SpecAgent<SpecType>;
    setAgent(newAgent: SpecAgent<SpecType>): void;
    setTalents(newTalents: SpecTalents<SpecType>): void;
    getTalentsString(): string;
    setTalentsString(newTalentsString: string): void;
    getSpecOptions(): SpecOptions<SpecType>;
    setSpecOptions(newSpecOptions: SpecOptions<SpecType>): void;
    lookupItemSpec(itemSpec: ItemSpec): EquippedItem | null;
    lookupEquipmentSpec(equipSpec: EquipmentSpec): Gear;
    computeGemEP(gem: Gem): number;
    computeEnchantEP(enchant: Enchant): number;
    computeItemEP(item: Item): number;
    makeCurrentIndividualSimRequest(iterations: number, debug: boolean): IndividualSimRequest;
    setWowheadData(equippedItem: EquippedItem, elem: HTMLElement): void;
    toJson(): Object;
    fromJson(obj: any): void;
}
