import React, {forwardRef, ReactNode, useCallback, useEffect, useImperativeHandle, useMemo, useState} from "react";
import './device.scss';
import ICON, {getIconByName} from "../../../icons/icons";
import IDevice, {Supplemental} from "../../../model/device";
import {useServices} from "../../../provider/service-provider";
import useEvents from "../../../hook/use-events";
import {AppEvent, EventTopic} from "../../../service/event-service";
import TimeoutIndicator from "../../timeout-indicator/timeout-indicator";
import useDateTimeFormat from "../../../hook/use-date-format";
import {Collapse} from "@material-ui/core";
import "./device.scss"
import {delay, first} from "rxjs/operators";
import useWebsocket from "../../../hook/use-websocket";
import useTranslator from "../../../hook/use-translator";
import {noop} from "rxjs";
import {sprintf} from "sprintf-js";
import Utils from "../../../utils/utils";
import RenderService from "../../../service/render-service";

function getDeviceIconClass(device: IDevice, timeout: boolean, deviceState: any): string {
    if (timeout || !deviceState) {
        return '';
    }

    let isOn = deviceState?.active;
    if (!!device.invert?.["state_color"]) {
        isOn = !isOn;
    }
    if (isOn) {
        return ' device-on';
    }
    return '';
}

function formatValue(value: any, supplemental: Supplemental, renderService: RenderService) {
    if (supplemental.format && supplemental.format.length && value != null) {
        return sprintf(supplemental.format, value)
    }
    if (supplemental.renderer && supplemental.renderer.length && value != null) {
        return renderService.render(supplemental.renderer, value);
    }

    return value;
}

export interface IDeviceControl {
    toggleDevice(active?: boolean): void;
}

interface DeviceProperties {
    device: IDevice;
    deviceToggle?: { getIcon: (timeout: boolean, devState: any) => React.ReactNode, onToggle?: () => void }
    expanderChildren?: ReactNode;
    children?: ReactNode;
    timeoutIndicator?: { disabled?: boolean, x?: number, y?: number };
    onTimeout?: (value: boolean) => void;
    onState?: (state: any) => void;
    onServerState?: (state: any) => any;
    toggleConfirm?: { title: string, message?: { activate?: string, deactivate?: string } },
}

const DeviceControl = forwardRef<IDeviceControl, DeviceProperties>((props, ref) => {
    const {
        device,
        expanderChildren,
        children,
        timeoutIndicator,
        onState,
        onServerState,
        onTimeout,
        deviceToggle,
        toggleConfirm
    } = props;
    const services = useServices();
    const formatDate = useDateTimeFormat();
    const websocket = useWebsocket();
    const translate = useTranslator();
    const [deviceState, setDeviceState] = useState(null);
    const [timeout, setTimeout] = useState(false);
    const [expanded, setExpanded] = useState(!!device.expanded);

    const setServerState = useMemo(() => (state: any) => {
        if (state && onServerState) {
            state = onServerState(state)
        }
        setDeviceState(() => state || {});
        onState && onState(state);
    }, [onServerState, onState]);

    const handleExpand = useCallback((value: boolean) => {
        setExpanded(value);
        services.user().saveUserSettings({
            settings: [{key: device.deviceKey + '.expanded', value: value ? '1' : '0'}]
        });
    }, [services, device])

    useEvents(() => [
        {
            topic: EventTopic.DeviceState,
            handler: (e: AppEvent) => {
                if (device.deviceKey === e.payload?.deviceKey) {
                    setServerState(e.payload?.state);
                }
            }
        },
    ], [onTimeout, timeoutIndicator, setServerState]);

    const sendToggleDevice = useMemo(() => (active?: boolean) => {
        const newActiveState = active != null ? active : !(deviceState?.active);
        let nextFunc = (value: boolean) => {
            if (value && device.optimistic) {
                const state = Object.assign({}, (deviceState || {}), {active: newActiveState});
                setDeviceState(() => state);
                onState && onState(state);
            }
        }
        if (!websocket.isConnected() && !device.optimistic) {
            nextFunc = (_: boolean) =>
                // why do i need a delay here ?
                services.device().getDeviceState(device.deviceKey).pipe(delay(2000), first())
                    .subscribe((state) => setServerState(state));
        }
        services.device().setDevicePower(device.deviceKey, newActiveState).pipe(first()).subscribe(nextFunc);
    }, [services, device, deviceState, onState, setServerState, websocket]);

    const toggleDevice = useCallback((active?: boolean) => {
        if (deviceToggle && deviceToggle.onToggle) {
            deviceToggle.onToggle();
            return;
        }
        if (deviceState) {
            const newState = active != null ? active : !deviceState.active;
            if (device.confirm) {
                const title = toggleConfirm?.title || 'text.default-toggle-tittle';
                const activateMsg = toggleConfirm?.message?.activate || 'text.turn-on-question';
                const deactivateMsg = toggleConfirm?.message?.deactivate || 'text.turn-off-question';
                services.dialog().confirm({
                    title: translate(title),
                    message: translate(newState ? activateMsg : deactivateMsg),
                    cancelLabel: translate('button.no'),
                    successLabel: translate('button.yes'),
                    onSubmit: () => sendToggleDevice(newState)
                })
            } else {
                sendToggleDevice(newState);
            }
        }
    }, [services, device, deviceState, deviceToggle, toggleConfirm, sendToggleDevice, translate]);

    useImperativeHandle(ref, () => ({toggleDevice}));

    useEffect(() => {
        if (timeoutIndicator?.disabled !== true) {
            const result = services.device().isDeviceTimeout(device, deviceState);
            setTimeout(() => result);
            onTimeout && onTimeout(result);
        }
        return noop
    }, [device, timeoutIndicator, deviceState, onTimeout, services])

    useEffect(() => {
        services.user().getUserSetting(device.deviceKey + '.expanded').pipe(first()).subscribe((value) => {
            if (!Utils.isBlank(value)) {
                setExpanded(parseInt(value) !== 0);
            }
        })
    }, [services, device.deviceKey]);

    const handleToggleClick = useCallback(() => {
        toggleDevice();
    }, [toggleDevice]);

    const handleExpanderClick = useCallback(() => {
        handleExpand(!expanded)
    }, [handleExpand, expanded]);

    let toggle: React.ReactNode = null;
    if (deviceToggle) {
        toggle = <div className={'device' + getDeviceIconClass(device, timeout, deviceState)}
                      onClick={handleToggleClick}>{getIconByName(device.icon) || deviceToggle.getIcon(timeout, deviceState)}</div>
    }
    return <div className={'device-panel'}>
        {toggle}
        {children}
        {timeout ? <TimeoutIndicator shiftX={timeoutIndicator?.x || 0} shiftY={timeoutIndicator?.y || 0}/> : null}
        <div className={'expander'}
             onClick={handleExpanderClick}>{expanded ? ICON.ArrowUp : ICON.ArrowDown}</div>
        <Collapse in={expanded} timeout="auto" unmountOnExit>
            <div className={'device-state-panel'}>
                <div className={'last-updated'}>{formatDate(deviceState?.lastUpdated, true)}</div>
                {expanderChildren}
                {device.supplemental?.length && device.supplemental.map(s => {
                    const v = deviceState?.[s.field];
                    if (v != null) {
                        return <div key={device.deviceKey + '-' + s.field} className={'device-state'}>
                            <div className={'device-state-label'}>{translate(s.caption)}</div>
                            <div className={'device-state-value'}>{formatValue(v, s, services.render())}</div>
                        </div>
                    }
                    return null;
                })}
            </div>
        </Collapse>
    </div>;
});

export default DeviceControl;
