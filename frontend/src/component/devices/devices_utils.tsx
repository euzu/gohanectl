import IRoom, {IconConfig, IRoomConfig} from "../../model/room";
import React from "react";
import {getIconByName} from "../../icons/icons";
import IDeviceList from "../../model/device-list";
import Utils from "../../utils/utils";
import EventService, {EventTopic} from "../../service/event-service";

const labelUnsorted = 'label.unsorted-room';
const labelAll = 'label.all-room';
const iconUnsorted = 'Unsorted';
const iconAll = 'All';

const swipeOptions = {
    delta: 80,                          // min distance(px) before a swipe starts
    preventDefaultTouchmoveEvent: true,  // call e.preventDefault *See Details*
    trackTouch: true,                    // track touch input
    trackMouse: true,                    // track mouse input
    rotationAngle: 0,                    // set a rotation angle
}

function getIcon(iconConfig: IconConfig): React.ReactElement {
    if (iconConfig) {
        if (iconConfig.key) {
            const icon = getIconByName(iconConfig.key);
            if (icon) {
                const myStyle = iconConfig.color ? {color: iconConfig.color} : {};
                return <div className={"room-icon"} style={myStyle}>{icon}</div>
            }
        }
    }
    return null;
}

const createRooms = (roomConfig: IRoomConfig, devCfg: IDeviceList): IRoom[] => {
    const result: IRoom[] = [];
    if (roomConfig && devCfg && roomConfig.definition && devCfg.devices) {
        const availableRooms: { [key: string]: boolean } = {[labelAll]: true};
        for (let i = 0, cnt = devCfg.devices.length; i < cnt; i++) {
            const d = devCfg.devices[i];
            if (Utils.isBlank(d.room)) {
                availableRooms[labelUnsorted] = true;
            } else {
                availableRooms[d.room] = true;
            }
        }
        const defaultIconOnly = !!roomConfig?.icon_only;
        for (let i = 0, cnt = roomConfig.definition.length; i < cnt; i++) {
            const cfg = roomConfig.definition[i];
            if (cfg.disabled !== true) {
                let caption = cfg.label;
                if (caption === iconAll) {
                    caption = labelAll;
                } else if (caption === iconUnsorted) {
                    caption = labelUnsorted;
                }
                if (availableRooms[caption]) {
                    let iconOnly = defaultIconOnly;
                    const room = {
                        key: cfg.key,
                        caption,
                        iconOnly,
                        icon: getIcon(cfg)
                    } as IRoom;
                    result.push(room);
                }
            }
        }
    }
    return result;
}

function createTabsetDef(rooms: IRoom[]) {
    if (rooms) {
        return rooms.map((room) => {
            return {
                icon: room.icon,
                label: room.iconOnly ? '' : room.caption,
            }
        })
    }
    return null;
}

function getVisibilities(tabIndex: number, devicesConfig: IDeviceList, rooms: IRoom[]) {
    const room = rooms[tabIndex];
    const all = room.caption === DevicesUtils.labelAll;
    const unsorted = !all && room.caption === DevicesUtils.labelUnsorted;
    const vis :{[key: string] : boolean} = {};
    for (let i = 0, cnt = devicesConfig?.devices?.length; i < cnt; i++) {
        const dev = devicesConfig.devices[i];
        let val = false;
        if (!all) {
            if (!Utils.isBlank(dev.room)) {
                val = dev.room === room.caption;
            } else if (unsorted) {
                val = true;
            }
        } else {
            val = true;
        }
        vis[dev.deviceKey] = val;
    }
    return vis;
}


function fireDeviceStates(deviceStates: any, eventService: EventService)  {
    for (let deviceKey in deviceStates) {
        if (deviceStates.hasOwnProperty(deviceKey)) {
            const state = deviceStates[deviceKey];
            if (state) {
                setTimeout(() =>
                    eventService.fireEvent(EventTopic.DeviceState, {deviceKey, state}), 500);
            }
        }
    }
}

const DevicesUtils = {
    swipeOptions,
    labelUnsorted,
    labelAll,
    iconUnsorted,
    iconAll,
    getIcon,
    createRooms,
    createTabsetDef,
    fireDeviceStates,
    getVisibilities
}

export default DevicesUtils;