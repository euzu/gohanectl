import React from "react";
import './light.scss';
import ICON from "../../../icons/icons";
import IDevice from "../../../model/device";
import DeviceControl from "../device/device";

interface LightProperties {
    device: IDevice;
}

const timeoutIndicator = {x: -40, y: -35};
const toggleConfirm = {title: 'label.light'};
const deviceToggle = {
    getIcon: (timeout: boolean,) => timeout ? ICON.BulbUnknown : ICON.Bulb,
};

export default function Light(props: LightProperties) {
    const {device} = props;

    return <DeviceControl device={device}
                          timeoutIndicator={timeoutIndicator}
                          deviceToggle={deviceToggle}
                          toggleConfirm={toggleConfirm}
    />;
}
