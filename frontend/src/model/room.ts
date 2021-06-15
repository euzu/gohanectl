import IDevice from "./device";

export default interface IRoom {
    key: string;
    caption: string;
    iconOnly: boolean;
    icon: React.ReactElement;
}

export interface IconConfig {
    label: string;
    key: string;
    color: string;
    disabled?: boolean;
}

export interface IRoomConfig {
    icon_only: boolean
    definition: IconConfig[]
}