export enum DeviceType {
    light = "light",
    socket = "socket",
    thermostat = "thermostat",
    temperature = "temperature",
    sensor = "sensor",
    roller_shutter = "roller-shutter",
}

export interface Supplemental {
    field: string,
    caption: string,
    format?: string,
    renderer?: string
}

export default interface IDevice {
    type: DeviceType,
    deviceKey: string,
    caption: string
    confirm: boolean;
    optimistic: boolean;
    url: string;
    timeout: number;
    invert: { position: boolean, 'state_color': boolean };
    room: string;
    supplemental: Supplemental[]
    expanded: boolean;
    icon: string;
}
