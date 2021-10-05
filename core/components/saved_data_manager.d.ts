import { Spec } from '../proto/common.js';
import { Sim } from '../sim.js';
import { TypedEvent } from '../typed_event.js';
import { Component } from '../components/component.js';
export declare type SavedDataManagerConfig<SpecType extends Spec, T> = {
    label: string;
    changeEmitters: Array<TypedEvent<any>>;
    equals: (a: T, b: T) => boolean;
    getData: (sim: Sim<SpecType>) => T;
    setData: (sim: Sim<SpecType>, data: T) => void;
    toJson: (a: T) => any;
    fromJson: (obj: any) => T;
};
export declare class SavedDataManager<SpecType extends Spec, T> extends Component {
    private readonly sim;
    private readonly config;
    private readonly userData;
    private readonly presets;
    private readonly savedDataDiv;
    private readonly saveInput;
    private frozen;
    constructor(parent: HTMLElement, sim: Sim<Spec>, config: SavedDataManagerConfig<SpecType, T>);
    addSavedData(newName: string, data: T, isPreset: boolean, tooltipInfo?: string): void;
    private makeSavedData;
    private saveUserData;
    loadUserData(): void;
    freeze(): void;
}
