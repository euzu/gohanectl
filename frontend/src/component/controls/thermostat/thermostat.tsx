import React, {useEffect, useState} from "react";
import './thermostat.scss';
import IDevice from "../../../model/device";
import "@thomasloven/round-slider";
import ICON from "../../../icons/icons";
import {useServices} from "../../../provider/service-provider";
import {noop} from "rxjs";
import {first} from "rxjs/operators";
import DeviceControl from "../device/device";
import {EventTopic} from "../../../service/event-service";

const timeoutIndicator = {x: -60, y: -50};

interface ThermostatProperties {
    device: IDevice;
}

export default function Thermostat(props: ThermostatProperties) {
    const {device} = props;
    const services = useServices();
    const [deviceState, setDeviceState] = useState(null);
    const [deviceTimeout, setDeviceTimeout] = useState(false);
    const [desiredTemp, setDesiredTemp] = useState((deviceState?.desired || 5));

    useEffect(() => {
        const valueChanging = (event: any) => {
            setDesiredTemp(event.detail.value);
            services.event().fireEvent(EventTopic.PreventSwipe, true);
        }
        const valueChanged = (event: any) => {
            const desiredTemp = event.detail.value;
            services.device().setDeviceTemperature(device.deviceKey, desiredTemp).pipe(first()).subscribe(noop);
            setTimeout(() => services.event().fireEvent(EventTopic.PreventSwipe, false), 1000);
        };
        const thermostat = document.getElementById('thermostat-' + device.deviceKey);
        // @ts-ignore
        thermostat.addEventListener('value-changing', valueChanging);
        // @ts-ignore
        thermostat.addEventListener('value-changed', valueChanged);
        // @ts-ignore
        return () => {
            thermostat.removeEventListener('value-changing', valueChanging);
            thermostat.removeEventListener('value-changed', valueChanged);
        }
    }, [device.deviceKey, services])

    const setState = (state: any) => {
        setDeviceState(() => state);
        setDesiredTemp(state.desired);
    }

    return <DeviceControl onTimeout={setDeviceTimeout} onState={setState} device={device}
                          timeoutIndicator={timeoutIndicator}>
        <div className={"thermostat-panel" + (deviceTimeout ? ' thermostat-timeout' : '')}>
            <div className={"thermostat"}>
                {/*
            // @ts-ignore */}
                <round-slider
                    id={'thermostat-' + device.deviceKey}
                    value={desiredTemp}
                    min={5}
                    max={30}
                    step={0.5}
                    arcLength={180}
                    startAngle={180}
                    handleSize={8}
                />
                <div className={"thermostat-row thermostat-desired"}>
                    <div className={"thermostat-value"}>{desiredTemp?.toFixed(1) || ''}</div>
                    <div className={"thermostat-unit"}>°C</div>
                </div>
            </div>
            <div className={"thermostat-row thermostat-temp"}>
                <div className={"thermostat-icon"}>{ICON.Thermometer}</div>
                <div className={"thermostat-value"}>{deviceState?.temp}</div>
                <div className={"thermostat-unit"}>°C</div>
            </div>

            <div className={"thermostat-row thermostat-valve"}>
                <div className={"thermostat-icon"}>{ICON.Thermostat}</div>
                <div className={"thermostat-value"}>{deviceState?.valve}</div>
            </div>
        </div>
    </DeviceControl>;
}
