import React, {useState} from "react";
import './socket.scss';
import ICON from "../../../icons/icons";
import IDevice from "../../../model/device";
import DeviceControl from "../device/device";
import useTranslator from "../../../hook/use-translator";

interface SocketProperties {
    device: IDevice;
}

const PROPERTIES = [
    {caption: 'label.voltage', field: 'voltage', unit: 'V'},
    {caption: 'label.current', field: 'current', unit: 'A'},
    {caption: 'label.power', field: 'power', unit: 'W'},
    {caption: 'label.apparent-power', field: 'apparentpower', unit: 'VA'},
    {caption: 'label.reactive-power', field: 'reactivepower', unit: 'VAr'},
    {caption: 'label.power-factor', field: 'factor'},
    {caption: 'label.yesterday', field: 'yesterday', unit: 'kWh'},
    {caption: 'label.today', field: 'today', unit: 'kWh'},
    {caption: 'label.total', field: 'total', unit: 'kWh'},
]

const timeoutIndicator = {x: -50, y: -35};
const toggleConfirm = {title: 'label.socket'};
const deviceToggle = {
    getIcon: (timeout: boolean,) => timeout ? ICON.SocketUnknown : ICON.Socket,
};

export default function Socket(props: SocketProperties) {
    const {device} = props;
    const translate = useTranslator();
    const [deviceState, setDeviceState] = useState(null);

    const expanderChildren = <React.Fragment>
        {PROPERTIES.map(prop => {
            const val = deviceState?.[prop.field];
            if (val != null) {
                return <div key={device.deviceKey + '-' + prop.field} className={'device-state'}>
                    <div className={'device-state-label'}>{translate(prop.caption)}</div>
                    <div className={'device-state-value'}>{val} {prop.unit}</div>
                </div>
            }
            return null;
        }
        ).filter(x => x!=null)}
    </React.Fragment>;

    return <DeviceControl onState={(s, )  => setDeviceState(s)} device={device}
                          expanderChildren={expanderChildren}
                          timeoutIndicator={timeoutIndicator}
                          deviceToggle={deviceToggle}
                          toggleConfirm={toggleConfirm}
    />

}
