import React, {useMemo, useState} from "react";
import './temperature.scss';
import DeviceControl from "../device/device";
import IDevice from "../../../model/device";
import {sprintf} from "sprintf-js";
import ICON from "../../../icons/icons";

const timeoutIndicator = {x: -50, y: -35};
const toggleConfirm = {title: 'label.socket'};

function formatValue(value: any) {
    if (value != null) {
        return sprintf("%.1f", value)
    }
    return value;
}

interface TemperatureProperties {
    device: IDevice;
}

export default function Temperature(props: TemperatureProperties) {
    const {device} = props;
    const [, setDeviceState] = useState(null);

    const deviceToggle = useMemo(() => {
        return {
            getIcon: (timeout: boolean, devState: any): React.ReactNode => {
                if (devState) {
                    return <div className={'temperature-display'}>{formatValue(devState[device.icon])}</div>;
                }
                return ICON.Thermometer
            }
        }
    }, [device.icon]);

    return <DeviceControl onState={(s,) => setDeviceState(s)}
                          device={device}
                          timeoutIndicator={timeoutIndicator}
                          deviceToggle={deviceToggle}
                          toggleConfirm={toggleConfirm}
    />
}
