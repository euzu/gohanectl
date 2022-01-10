import React, {useEffect, useMemo, useRef, useState} from "react";
import "./devices.scss"
import {useServices} from "../../provider/service-provider";
import IDeviceList from "../../model/device-list";
import IRoom, {IRoomConfig} from '../../model/room';
import {AppEvent, EventTopic} from "../../service/event-service";
import {first} from "rxjs/operators";
import useWakeup from "../../hook/use-wakeup";
import {combineLatest, noop} from "rxjs";
import TabSet, {ITabSet} from "../tab-set/tab-set";
import {useSwipeable} from "react-swipeable";
import DeviceCard from "../device-card/device-card";
import IDevice from "../../model/device";
import useEvents from "../../hook/use-events";
import DevicesUtils from "./devices_utils";

const TAB_INDEX_KEY: string = 'hanectl-current-tab';

function storeTabIndex(index: number): void {
    localStorage.setItem(TAB_INDEX_KEY, ''+index);
}

function loadTabIndex(): number {
    return parseInt(localStorage.getItem(TAB_INDEX_KEY)) || 0;
}

interface DevicesProperties {

}

export default function Devices(params: DevicesProperties) {
    const services = useServices();
    const [authenticated, setAuthenticated] = useState<boolean>(false);
    const [devicesConfig, setDevicesConfig] = useState<IDeviceList>(null);
    const [rooms, setRooms] = useState<IRoom[]>(null);
    const [tabIndex, setTabIndex] = useState<number>(loadTabIndex());
    const [tabs, setTabs] = useState<React.ReactNode>(null);
    const [visibility, setVisibility] = useState<{[key: string] : boolean}>({});
    const wakeup = useWakeup();
    const tabsetRef = useRef<ITabSet>();
    const preventSwipe = useRef<boolean>(false);

    const changeTabIndex = useMemo(() => (index: number) => {
        setTabIndex(index);
        storeTabIndex(index);
    }, []);

    const setRoomsConfig = useMemo(() => (roomsCfg: IRoomConfig, devCfg: IDeviceList) => {
        const newRooms = DevicesUtils.createRooms(roomsCfg, devCfg);
        const tabsDef = DevicesUtils.createTabsetDef(newRooms);
        let roomTabs: React.ReactNode;
        if (tabsDef && tabsDef.length) {
            roomTabs =
                <TabSet ref={tabsetRef} index={tabIndex} tabs={tabsDef} rolling={true} onChange={changeTabIndex}/>
        }
        setRooms(() => newRooms)
        setTabs(() => roomTabs);
    }, [changeTabIndex, tabIndex]);


    useEffect(() => {
        const sub = services.auth().isAuthenticated().subscribe(val => setAuthenticated(val));
        return () => sub.unsubscribe();
    }, [authenticated, services])

    useEffect(() => {
        if (authenticated) {
            combineLatest([
                services.device().getDeviceConfig().pipe(first()),
                services.config().getRoomConfig().pipe(first()),
                services.device().getDeviceStates().pipe(first()),
            ]).pipe(first()).subscribe({
                next: (data) => {
                    const [dev, room, states] = data;
                    setDevicesConfig(() => dev);
                    setRoomsConfig(room, dev);
                    DevicesUtils.fireDeviceStates(states, services.event());

                },
                error: noop
            });
        }
        return noop;
    }, [authenticated, services, setRoomsConfig])

    useEffect(() => {
        if (devicesConfig && rooms && rooms.length && devicesConfig.devices?.length) {
            const vis = DevicesUtils.getVisibilities(tabIndex, devicesConfig, rooms);
            setVisibility(vis);
        }
        return noop;
    }, [tabIndex, devicesConfig, rooms])

    useEffect(() => {
        const subs = wakeup.subscribe(() => {
            if (authenticated) {
                services.device().getDeviceStates().pipe(first()).subscribe({
                    next: (states) => DevicesUtils.fireDeviceStates(states, services.event()),
                    error: noop
                })
            }
        });
        return () => subs.unsubscribe();
    }, [authenticated, wakeup, services]);

    useEvents(() => [
        {
            topic: EventTopic.PreventSwipe,
            handler: (e: AppEvent) => preventSwipe.current = !!e.payload
        }
    ], [preventSwipe]);


    const handlers = useSwipeable({
        onSwipedLeft: () => !preventSwipe.current && tabsetRef.current?.setNextTab(),
        onSwipedRight: () => !preventSwipe.current && tabsetRef.current?.setPrevTab(),
        ...DevicesUtils.swipeOptions,
    });

    return <div className={"devices-container"} {...handlers}>
        <div className={"tabs-panel"}>{tabs}</div>
        <div className={"devices-panel"} {...handlers}>
            {devicesConfig?.devices?.map((device: IDevice) => <DeviceCard key={'card-' + device.deviceKey}
                                                                          hidden={!visibility[device.deviceKey]} device={device}/>)}
        </div>
    </div>
}
