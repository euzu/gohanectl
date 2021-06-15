import React, {useMemo} from "react";
import './device-card.scss';
import IDevice, {DeviceType} from "../../model/device";
import Light from "../controls/light/light";
import Thermostat from "../controls/thermostat/thermostat";
import Socket from "../controls/socket/socket";
import Sensor from "../controls/sensor/sensor";
import ICON from "../../icons/icons";
import RollerShutter from "../controls/roller-shutter/roller-shutter";
import Temperature from "../controls/temperature/temperature";

interface DeviceCardProperties {
    device: IDevice;
    hidden: boolean;
}

export default function DeviceCard(props: DeviceCardProperties) {
    const {device, hidden} = props;
    const createDevice = useMemo(() => (device: IDevice) => {
        switch (device.type) {
            case DeviceType.light:
                return <Light key={'light-' + device.deviceKey} device={device}/>
            case DeviceType.socket:
                return <Socket key={'socket-' + device.deviceKey} device={device}/>
            case DeviceType.thermostat:
                return <Thermostat key={'thermostat-' + device.deviceKey} device={device}/>
            case DeviceType.temperature:
                return <Temperature key={'temperature-' + device.deviceKey} device={device}/>
            case DeviceType.sensor:
                return <Sensor key={'sensor-' + device.deviceKey} device={device}/>
            case DeviceType.roller_shutter:
                return <RollerShutter key={'sensor-' + device.deviceKey} device={device}/>
        }
        return null
    }, []);

    const renderCaption = () => {
        if (device.url) {
            return <React.Fragment>
                <a href={device.url} target={'_blank'} rel={'noreferrer'}>{device.caption}</a>
                <div className={"device-caption-icon"}>{ICON.Link}</div>
            </React.Fragment>
        }
        return device.caption
    }

    return <div className={"device-card" + (hidden ? ' hidden' : '')}>
        <div className={"device-caption"}>{renderCaption()}</div>
        <div className={"device-card-content"}>
            {createDevice(device)}
        </div>
    </div>;
}
