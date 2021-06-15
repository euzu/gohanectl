import React, {useEffect, useState} from 'react';
import './header.scss';
import {useServices} from "../../provider/service-provider";
import useEvents from "../../hook/use-events";
import {AppEvent, EventTopic} from "../../service/event-service";
import ICON from "../../icons/icons";
import {first} from "rxjs/operators";
import {noop} from "rxjs";
import ServerState from "../../model/server-state";

export interface HeaderProperties {
}

export default function Header(props: HeaderProperties) {
    const services = useServices();
    const [pingState, setPingState] = useState(false);
    const [mqttState, setMqttState] = useState(1);

    useEffect(() => {
        const subs = services.auth().isAuthenticated().subscribe((val: boolean) => {
            if (val) {
                services.config().getServerState().pipe(first()).subscribe({
                    next: (data: ServerState) => {
                        setPingState(data.websocket);
                        setMqttState(data.mqtt);
                    },
                    error: noop
                });
            }
        });
        return () => subs.unsubscribe();
    }, [services])

    useEvents(() => [
        {
            topic: EventTopic.Websocket,
            handler: (e: AppEvent) => {
                if (e.payload) {
                    const data = e.payload;
                    if (data.error || data.connected === false) {
                        setPingState(() => false);
                    } else if (data.ping || data.connected === true) {
                        setPingState(() => true);
                    }
                }
            }
        },
        {
            topic: EventTopic.Mqtt,
            handler: (e: AppEvent) => setMqttState(e.payload)
        },
    ], [services]);

    return <div className={'header'}>
        <div className={'header-logo'}>
            <div className={'header-logo-icon'}>{ICON.Logo}</div>
            <div className={'header-logo-text'}>Hane-Control</div>
        </div>
        <div className={'header-tools'}>
            <div className={'header-icon ' + (mqttState === 1 ? 'mqtt-on' : 'mqtt-off')}>{ICON.Mqtt}</div>
            <div className={'header-icon ' + (pingState ? 'ping-on' : 'ping-off')}>{ICON.Ping}</div>
        </div>
    </div>

}
