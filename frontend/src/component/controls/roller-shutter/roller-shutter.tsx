import React, {useMemo, useRef, useState} from "react";
import './roller-shutter.scss';
import IDevice from "../../../model/device";
import "@thomasloven/round-slider";
import ICON from "../../../icons/icons";
import {useServices} from "../../../provider/service-provider";
import {noop} from "rxjs";
import {first} from "rxjs/operators";
import DeviceControl, {IDeviceControl} from "../device/device";
import useTranslator from "../../../hook/use-translator";
import {Slider} from "@material-ui/core";
import {EventTopic} from "../../../service/event-service";

const MAX_POSITION = 100;

function clamp(value: number, min: number, max: number): number {
    return Math.max(min, Math.min(max, value));
}

const getIcon = (timeout: boolean, devState: any) => {
    if (timeout) {
        return ICON.RollerShutterUnknown;
    }
    let pos = devState?.position;
    if (typeof pos === "number") {
        if (pos < 25) {
            return ICON.RollerShutterClosed;
        }
        if (pos < 50) {
            return ICON.RollerShutter25;
        }
        if (pos < 75) {
            return ICON.RollerShutter50;
        }
        if (pos < 100) {
            return ICON.RollerShutter75;
        } else {
            return ICON.RollerShutterOpen;
        }
    }
    return devState?.active ? ICON.RollerShutterOpen : ICON.RollerShutterClosed
}

const timeoutIndicator = {x: -50, y: -90};
const deviceToggle = {getIcon}
const toggleConfirm = {
    title: 'label.roller-shutter',
    message: {
        activate: 'text.open-question',
        deactivate: 'text.close-question'

    }
}

interface RollerShutterProperties {
    device: IDevice;
}

export default function RollerShutter(props: RollerShutterProperties) {
    const {device} = props;
    const services = useServices();
    const translate = useTranslator();
    const deviceRef = useRef<IDeviceControl>();
    const [position, setPosition] = useState(0);
    const [sliderPos, setSliderPos] = useState(0);

    const setServerState = useMemo(() => (state: any): any => {
        const invert = device?.invert?.position;
        let newState: any = {position: 0};
        if (state) {
            newState = Object.assign(newState, state)
            if (state.position != null) {
                newState.position = invert ? (MAX_POSITION - newState.position) : newState.position;
            } else {
                newState.position = state.active ? MAX_POSITION : 0;
            }
        }
        newState.position = clamp(newState.position, 0, MAX_POSITION);
        newState.active = newState.position > 0;
        setPosition(() => newState.position);
        return newState;
    }, [device]);


    const setState = useMemo(() => (state: any) => {
        setSliderPos(() => state.position);
    }, []);

    const pauseDevice = () => {
        services.device().setDevicePause(device.deviceKey).pipe(first()).subscribe(noop);
    }

    const pauseRollerShutter = () => {
        if (device.confirm) {
            services.dialog().confirm({
                title: translate('label.roller-shutter'),
                message: translate('text.pause-question'),
                cancelLabel: translate('button.no'),
                successLabel: translate('button.yes'),
                onSubmit: () => pauseDevice()
            })
        } else {
            pauseDevice();
        }
    }

    const setSliderPosition = (value: number | number[]) => {
        services.event().fireEvent(EventTopic.PreventSwipe, true);
        if (typeof value === "number" && !isNaN(value)) {
            setSliderPos(value);
        }
    }

    const setShutterPosition = useMemo(() => (value: number) => {
        if (!isNaN(value)) {
            const positionDevice = () => {
                const position = device?.invert?.position ? (MAX_POSITION - value) : value;
                services.device().setDevicePosition(device.deviceKey, position).pipe(first()).subscribe(noop);
            }
            if (device.confirm) {
                services.dialog().confirm({
                    title: translate('label.roller-shutter'),
                    message: translate('text.position-question'),
                    cancelLabel: translate('button.no'),
                    successLabel: translate('button.yes'),
                    onSubmit: () => positionDevice(),
                    onCancel: () => setSliderPos(() => position)
                })
            } else {
                positionDevice();
            }
        }
        setTimeout(() => services.event().fireEvent(EventTopic.PreventSwipe, false), 1000);
    }, [device, services, translate, position])

    return <DeviceControl ref={deviceRef} onServerState={setServerState} onState={setState}
                          device={device}
                          timeoutIndicator={timeoutIndicator}
                          deviceToggle={deviceToggle}
                          toggleConfirm={toggleConfirm}
    >
        <div className={'roller-shutter-buttons'}>
            <div className={"roller-shutter-icon"} onClick={() => deviceRef.current.toggleDevice(false)}>{ICON.Up}</div>
            <div className={"roller-shutter-icon"} onClick={() => pauseRollerShutter()}>{ICON.Pause}</div>
            <div className={"roller-shutter-icon"}
                 onClick={() => deviceRef.current.toggleDevice(true)}>{ICON.Down}</div>
        </div>
        <div className={'roller-shutter-slider'}>
            <Slider
                min={0} max={MAX_POSITION} step={1} value={sliderPos}
                onMouseUp={() => setShutterPosition(sliderPos)}
                onTouchEnd={() => setShutterPosition(sliderPos)}
                onChange={(event, value) => setSliderPosition(value)}
                valueLabelDisplay="auto"
                aria-labelledby="roller-shutter-slider"/>
        </div>

    </DeviceControl>;
}
