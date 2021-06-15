import React from "react";
import './sensor.scss';
import ICON from "../../../icons/icons";
import IDevice from "../../../model/device";
import DeviceControl from "../device/device";

const timeoutIndicator = {disabled: true};
const deviceToggle = {
    getIcon: () => ICON.Sensor,
};

interface SensorProperties {
    device: IDevice;
}

export default function Sensor(props: SensorProperties) {
    const {device} = props;

    return <DeviceControl device={device}
                          timeoutIndicator={timeoutIndicator}
                          deviceToggle={deviceToggle}
    />
}
